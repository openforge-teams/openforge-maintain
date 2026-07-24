package file

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"
)

// FileItem represents a file or directory entry.
type FileItem struct {
	Name    string      `json:"name"`
	Path    string      `json:"path"`
	IsDir   bool        `json:"is_dir"`
	Size    int64       `json:"size"`
	ModTime time.Time    `json:"mod_time"`
	Mode    os.FileMode `json:"mode"`
	Owner   string      `json:"owner"`
}

// FileInfo holds detailed file information.
type FileInfo struct {
	FileItem
	IsSymlink bool   `json:"is_symlink"`
	MimeType  string `json:"mime_type"`
}

// FileManagerService provides file management operations.
type FileManagerService struct {
	baseDir string
}

// NewFileManagerService creates a new FileManagerService.
func NewFileManagerService(baseDir string) *FileManagerService {
	return &FileManagerService{baseDir: baseDir}
}

// ListDir lists the contents of a directory.
func (s *FileManagerService) ListDir(path string) ([]FileItem, error) {
	absPath := s.resolvePath(path)
	entries, err := os.ReadDir(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	items := make([]FileItem, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		uid, gid := getFileOwner(info)
		items = append(items, FileItem{
			Name:    entry.Name(),
			Path:    filepath.Join(path, entry.Name()),
			IsDir:   entry.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
			Mode:    info.Mode(),
			Owner:   fmt.Sprintf("%d:%d", uid, gid),
		})
	}

	// Sort: directories first, then by name
	sort.Slice(items, func(i, j int) bool {
		if items[i].IsDir != items[j].IsDir {
			return items[i].IsDir
		}
		return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
	})

	return items, nil
}

// GetFileContent reads and returns the content of a file.
func (s *FileManagerService) GetFileContent(path string) (string, error) {
	absPath := s.resolvePath(path)
	data, err := os.ReadFile(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(data), nil
}

// SaveFileContent writes content to a file.
func (s *FileManagerService) SaveFileContent(path, content string) error {
	absPath := s.resolvePath(path)

	dir := filepath.Dir(absPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return os.WriteFile(absPath, []byte(content), 0644)
}

// Upload saves uploaded content to a destination path.
func (s *FileManagerService) Upload(src io.Reader, destPath string, overwrite bool) error {
	absPath := s.resolvePath(destPath)

	if !overwrite {
		if _, err := os.Stat(absPath); err == nil {
			return fmt.Errorf("file already exists: %s", destPath)
		}
	}

	dir := filepath.Dir(absPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	dst, err := os.Create(absPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

// Download returns a ReadCloser for the specified file.
func (s *FileManagerService) Download(path string) (io.ReadCloser, error) {
	absPath := s.resolvePath(path)
	file, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	return file, nil
}

// Delete removes a file or directory.
func (s *FileManagerService) Delete(path string) error {
	absPath := s.resolvePath(path)
	return os.RemoveAll(absPath)
}

// Rename moves or renames a file or directory.
func (s *FileManagerService) Rename(oldPath, newPath string) error {
	absOld := s.resolvePath(oldPath)
	absNew := s.resolvePath(newPath)

	dir := filepath.Dir(absNew)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return os.Rename(absOld, absNew)
}

// Chmod changes the permissions of a file or directory.
func (s *FileManagerService) Chmod(path string, mode uint32) error {
	absPath := s.resolvePath(path)
	return os.Chmod(absPath, os.FileMode(mode))
}

// Chown changes the owner of a file or directory.
func (s *FileManagerService) Chown(path string, uid, gid int) error {
	absPath := s.resolvePath(path)
	return os.Chown(absPath, uid, gid)
}

// Compress creates an archive from the given paths.
func (s *FileManagerService) Compress(paths []string, dest, format string) error {
	absDest := s.resolvePath(dest)

	dir := filepath.Dir(absDest)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	destFile, err := os.Create(absDest)
	if err != nil {
		return fmt.Errorf("failed to create archive: %w", err)
	}
	defer destFile.Close()

	var writer io.Writer = destFile
	if format == "gz" || format == "tar.gz" || format == "tgz" {
		gzWriter := gzip.NewWriter(destFile)
		defer gzWriter.Close()
		writer = gzWriter
	}

	tarWriter := tar.NewWriter(writer)
	defer tarWriter.Close()

	for _, p := range paths {
		absPath := s.resolvePath(p)
		err := filepath.Walk(absPath, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			relPath, err := filepath.Rel(filepath.Dir(absPath), filePath)
			if err != nil {
				return err
			}

			header, err := tar.FileInfoHeader(info, relPath)
			if err != nil {
				return err
			}
			header.Name = filepath.Join(filepath.Base(p), relPath)

			if err := tarWriter.WriteHeader(header); err != nil {
				return err
			}

			if !info.IsDir() {
				file, err := os.Open(filePath)
				if err != nil {
					return err
				}
				defer file.Close()
				_, err = io.Copy(tarWriter, file)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("failed to add %s to archive: %w", p, err)
		}
	}

	return nil
}

// Extract extracts an archive to the specified destination.
func (s *FileManagerService) Extract(src, dest string) error {
	absSrc := s.resolvePath(src)
	absDest := s.resolvePath(dest)

	if err := os.MkdirAll(absDest, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	file, err := os.Open(absSrc)
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	var reader io.Reader = file
	if strings.HasSuffix(src, ".gz") || strings.HasSuffix(src, ".tgz") {
		gzReader, err := gzip.NewReader(file)
		if err != nil {
			return fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzReader.Close()
		reader = gzReader
	}

	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar entry: %w", err)
		}

		target := filepath.Join(absDest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		case tar.TypeReg:
			dir := filepath.Dir(target)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			outFile, err := os.Create(target)
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("failed to write file: %w", err)
			}
			outFile.Close()
		case tar.TypeSymlink:
			dir := filepath.Dir(target)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			os.Remove(target)
			if err := os.Symlink(header.Linkname, target); err != nil {
				return fmt.Errorf("failed to create symlink: %w", err)
			}
		}
	}

	return nil
}

// Mkdir creates a directory.
func (s *FileManagerService) Mkdir(path string) error {
	absPath := s.resolvePath(path)
	return os.MkdirAll(absPath, 0755)
}

// GetInfo returns detailed information about a file.
func (s *FileManagerService) GetInfo(path string) (*FileInfo, error) {
	absPath := s.resolvePath(path)
	info, err := os.Lstat(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	uid, gid := getFileOwner(info)
	return &FileInfo{
		FileItem: FileItem{
			Name:    info.Name(),
			Path:    path,
			IsDir:   info.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
			Mode:    info.Mode(),
			Owner:   fmt.Sprintf("%d:%d", uid, gid),
		},
		IsSymlink: (info.Mode() & os.ModeSymlink) != 0,
	}, nil
}

// resolvePath resolves a relative path to an absolute path within the base directory.
func (s *FileManagerService) resolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(s.baseDir, path)
}

// getFileOwner extracts the UID and GID from file info.
func getFileOwner(info os.FileInfo) (uint32, uint32) {
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		return stat.Uid, stat.Gid
	}
	return 0, 0
}
