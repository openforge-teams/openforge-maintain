import http from 'http';
import { readFileSync, existsSync } from 'fs';
import { join, extname } from 'path';
import { randomUUID } from 'crypto';

const FRONTEND_DIR = process.env.FRONTEND_DIR || '/workspace/openforge-maintain/frontend/dist';
const PORT = process.env.PORT || 9998;
const tokens = new Map();

function jsonRes(res, code, data) {
  res.writeHead(200, { 'Content-Type': 'application/json', 'Access-Control-Allow-Origin': '*', 'Access-Control-Allow-Headers': 'Content-Type,Authorization' });
  res.end(JSON.stringify({ code: 0, message: 'success', data, trace_id: randomUUID() }));
}
function errRes(res, code, msg) {
  res.writeHead(code, { 'Content-Type': 'application/json', 'Access-Control-Allow-Origin': '*' });
  res.end(JSON.stringify({ code, message: msg, data: null, trace_id: randomUUID() }));
}
function createToken(userId, username) {
  const token = 'jwt_' + randomUUID().replace(/-/g, '') + '_' + Date.now();
  tokens.set(token, { userId, username });
  return token;
}
function verifyToken(auth) {
  if (!auth || !auth.startsWith('Bearer ')) return null;
  return tokens.get(auth.slice(7)) || null;
}

const server = http.createServer(async (req, res) => {
  const url = new URL(req.url, `http://${req.headers.host}`);
  const path = url.pathname;
  const method = req.method;

  if (method === 'OPTIONS') {
    res.writeHead(204, { 'Access-Control-Allow-Origin': '*', 'Access-Control-Allow-Methods': 'GET,POST,PUT,DELETE,OPTIONS', 'Access-Control-Allow-Headers': 'Content-Type,Authorization' });
    res.end(); return;
  }
  if (path === '/api/health') { jsonRes(res, 200, { status: 'healthy', version: '1.0.0', uptime: Math.floor(Date.now() / 1000) }); return; }

  if (path === '/api/v2/core/auth/login' && method === 'POST') {
    let body = '';
    for await (const chunk of req) body += chunk;
    const { username, password } = JSON.parse(body || '{}');
    const token = createToken(1, username);
    jsonRes(res, 200, { access_token: token, refresh_token: token, user: { id: 1, username, email: `${username}@openforge.local`, role_id: 1 } });
    return;
  }
  if (path === '/api/v2/core/auth/logout' && method === 'POST') { jsonRes(res, 200, null); return; }
  if (path === '/api/v2/core/auth/profile' && method === 'GET') {
    const auth = verifyToken(req.headers.authorization);
    if (!auth) return errRes(res, 401, 'unauthorized');
    jsonRes(res, 200, { id: auth.userId, username: auth.username, email: 'admin@openforge.local', role_id: 1 });
    return;
  }

  // Dashboard
  if (path === '/api/v2/dashboard/overview') { jsonRes(res, 200, { cpu_usage: 23.5, memory_usage: 61.2, memory_total: 15.6, memory_used: 9.55, disk_usage: 45.8, disk_total: 256.0, disk_used: 117.2, container_count: 24, container_running: 12, container_stopped: 12, network_in: 128.5, network_out: 42.3, uptime: 259200, hostname: 'openforge-server', os: 'Ubuntu 22.04 LTS' }); return; }
  if (path === '/api/v2/dashboard/cpu') { const now = Math.floor(Date.now()/1000); jsonRes(res, 200, { usage: 23.5, cores: 4, model: 'Intel Xeon E5-2680 v4', history: Array.from({length:60},(_,i)=>({time:String(now-60+i),value:18+(i%20)*0.5})) }); return; }
  if (path === '/api/v2/dashboard/memory') { const now = Math.floor(Date.now()/1000); jsonRes(res, 200, { total: 15.6, used: 9.55, free: 2.85, cached: 3.2, buffers: 0.8, history: Array.from({length:60},(_,i)=>({time:String(now-60+i),value:58+(i%15)*0.4})) }); return; }
  if (path === '/api/v2/dashboard/disk') { jsonRes(res, 200, { total: 256.0, used: 117.2, free: 138.8, partitions: [{ device: '/dev/sda1', mount: '/', total: 256.0, used: 117.2, free: 138.8, usage: 45.8 }, { device: '/dev/sda2', mount: '/home', total: 512.0, used: 89.4, free: 422.6, usage: 17.5 }, { device: '/dev/sdb1', mount: '/data', total: 1024.0, used: 623.8, free: 400.2, usage: 60.9 }] }); return; }
  if (path === '/api/v2/dashboard/network') { const now = Math.floor(Date.now()/1000); jsonRes(res, 200, { interfaces: [{ name: 'eth0', rx_bytes: 128.5*1024*1024, tx_bytes: 42.3*1024*1024, rx_speed: 15.2, tx_speed: 5.8 }], history: Array.from({length:60},(_,i)=>({time:String(now-60+i),rx:(80+(i%40))*1.5,tx:(30+(i%20))*1.2})) }); return; }
  if (path === '/api/v2/dashboard/processes') { jsonRes(res, 200, [{ pid: 1, name: 'systemd', user: 'root', cpu: 0.0, memory: 1.2, status: 'running' },{ pid: 1023, name: 'nginx', user: 'www-data', cpu: 0.5, memory: 2.8, status: 'running' },{ pid: 2045, name: 'docker', user: 'root', cpu: 1.2, memory: 15.3, status: 'running' },{ pid: 3098, name: 'mysql', user: 'mysql', cpu: 2.1, memory: 28.7, status: 'running' },{ pid: 4521, name: 'redis-server', user: 'redis', cpu: 0.3, memory: 4.5, status: 'running' },{ pid: 5678, name: 'node', user: 'deploy', cpu: 3.8, memory: 12.1, status: 'running' },{ pid: 7201, name: 'sshd', user: 'root', cpu: 0.1, memory: 0.8, status: 'running' },{ pid: 8934, name: 'openforge-core', user: 'root', cpu: 0.8, memory: 5.2, status: 'running' }]); return; }

  // Containers
  if (path === '/api/v2/containers') {
    jsonRes(res, 200, { list: [
      { id: 'a1b2c3d4', name: 'nginx-proxy', image: 'nginx:1.25', state: 'running', status: 'Up 3 days', ports: [{ ip: '0.0.0.0', private_port: 80, public_port: 80, type: 'tcp' },{ ip: '0.0.0.0', private_port: 443, public_port: 443, type: 'tcp' }], created_at: '2024-07-21T10:30:00Z', labels: {} },
      { id: 'e5f6g7h8', name: 'mysql-8', image: 'mysql:8.0', state: 'running', status: 'Up 3 days', ports: [{ ip: '0.0.0.0', private_port: 3306, public_port: 3306, type: 'tcp' }], created_at: '2024-07-21T10:30:00Z', labels: {} },
      { id: 'i9j0k1l2', name: 'redis-7', image: 'redis:7-alpine', state: 'running', status: 'Up 3 days', ports: [{ ip: '0.0.0.0', private_port: 6379, public_port: 6379, type: 'tcp' }], created_at: '2024-07-21T10:30:00Z', labels: {} },
      { id: 'm3n4o5p6', name: 'wordpress', image: 'wordpress:6.4', state: 'running', status: 'Up 2 days', ports: [{ ip: '0.0.0.0', private_port: 80, public_port: 8080, type: 'tcp' }], created_at: '2024-07-22T10:30:00Z', labels: {} },
      { id: 'q7r8s9t0', name: 'halo', image: 'halohub/halo:2.15', state: 'running', status: 'Up 1 day', ports: [{ ip: '0.0.0.0', private_port: 8090, public_port: 8090, type: 'tcp' }], created_at: '2024-07-23T10:30:00Z', labels: {} },
      { id: 'u1v2w3x4', name: 'minio', image: 'minio/minio:latest', state: 'running', status: 'Up 5 days', ports: [{ ip: '0.0.0.0', private_port: 9000, public_port: 9000, type: 'tcp' },{ ip: '0.0.0.0', private_port: 9001, public_port: 9001, type: 'tcp' }], created_at: '2024-07-19T10:30:00Z', labels: {} },
      { id: 'y5z6a7b8', name: 'gitea', image: 'gitea/gitea:1.21', state: 'running', status: 'Up 4 days', ports: [{ ip: '0.0.0.0', private_port: 3000, public_port: 3000, type: 'tcp' },{ ip: '0.0.0.0', private_port: 22, public_port: 2222, type: 'tcp' }], created_at: '2024-07-20T10:30:00Z', labels: {} },
      { id: 'c9d0e1f2', name: 'n8n', image: 'n8nio/n8n:latest', state: 'stopped', status: 'Exited (0) 2 hours ago', ports: [{ ip: '0.0.0.0', private_port: 5678, public_port: 5678, type: 'tcp' }], created_at: '2024-07-24T10:30:00Z', labels: {} },
    ], total: 8, page: 1, page_size: 20 });
    return;
  }

  // Static files
  if (path.startsWith('/assets/') || ['.js','.css','.png','.svg','.ico','.woff2','.woff','.json'].some(e => path.endsWith(e))) {
    try { const data = readFileSync(join(FRONTEND_DIR, path)); const ext = extname(path); const types = { '.html':'text/html','.js':'application/javascript','.css':'text/css','.png':'image/png','.svg':'image/svg+xml','.ico':'image/x-icon','.json':'application/json','.woff2':'font/woff2','.woff':'font/woff' }; res.writeHead(200, { 'Content-Type': types[ext] || 'application/octet-stream' }); res.end(data); return; } catch {}
  }

  // SPA fallback
  try { const data = readFileSync(join(FRONTEND_DIR, 'index.html')); res.writeHead(200, { 'Content-Type': 'text/html' }); res.end(data); } catch { res.writeHead(404); res.end('Not Found'); }
});

server.listen(PORT, () => console.log(`Server on :${PORT}, frontend: ${FRONTEND_DIR}`));
