import { useState, useEffect } from 'react';
import './App.css';
import { ConnectVPN, DisconnectVPN, GetVPNStatus, ActivateLicense, GetLicenseInfo } from '../wailsjs/go/main/App';

function App() {
  const [serverIP, setServerIP] = useState('');
  const [port, setPort] = useState('');
  const [password, setPassword] = useState('');
  const [isConnected, setIsConnected] = useState(false);
  const [stats, setStats] = useState({});
  const [licenseKey, setLicenseKey] = useState('');
  const [licenseInfo, setLicenseInfo] = useState({});
  const [error, setError] = useState('');

  useEffect(() => {
    // Check initial connection status
    updateStatus();
    // Check license status and info
    checkLicense();
    // Update stats periodically
    const interval = setInterval(updateStatus, 1000);
    return () => clearInterval(interval);
  }, []);

  const updateStatus = async () => {
    try {
      const status = await GetVPNStatus();
      setIsConnected(status.connected);
      if (status.connected && status.stats) {
        setStats(status.stats);
      }
    } catch (err) {
      setError(err.toString());
    }
  };

  const checkLicense = async () => {
    try {
      const info = await GetLicenseInfo();
      setLicenseInfo(info);
    } catch (err) {
      setError(err.toString());
    }
  };

  const handleConnect = async () => {
    try {
      setError('');
      const result = await ConnectVPN(serverIP, parseInt(port), password);
      if (result.success) {
        await updateStatus();
      } else {
        setError(result.error);
      }
    } catch (err) {
      setError(err.toString());
    }
  };

  const handleDisconnect = async () => {
    try {
      setError('');
      const result = await DisconnectVPN();
      if (result.success) {
        await updateStatus();
      } else {
        setError(result.error);
      }
    } catch (err) {
      setError(err.toString());
    }
  };

  const handleActivateLicense = async () => {
    try {
      setError('');
      const result = await ActivateLicense(licenseKey);
      if (result.success) {
        setLicenseInfo(result.info);
        setLicenseKey('');
      } else {
        setError(result.error);
      }
    } catch (err) {
      setError(err.toString());
    }
  };

  return (
    <div className="container">
      <h1>Outline VPN Client</h1>

      {/* License Section */}
      <div className="section">
        <h2>授權資訊</h2>
        {licenseInfo.status === '未授權' ? (
          <div className="license-activation">
            <input
              type="text"
              value={licenseKey}
              onChange={(e) => setLicenseKey(e.target.value)}
              placeholder="輸入授權碼"
            />
            <button onClick={handleActivateLicense}>啟動</button>
          </div>
        ) : (
          <div className="license-info">
            <p>狀態: {licenseInfo.status}</p>
            <p>使用者: {licenseInfo.issuedTo}</p>
            <p>到期時間: {licenseInfo.expiresIn}</p>
          </div>
        )}
      </div>

      {/* VPN Connection Section */}
      {licenseInfo.status === '已授權' && (
        <div className="section">
          <h2>VPN 連接</h2>
          {!isConnected ? (
            <div className="connection-form">
              <input
                type="text"
                value={serverIP}
                onChange={(e) => setServerIP(e.target.value)}
                placeholder="伺服器 IP"
              />
              <input
                type="text"
                value={port}
                onChange={(e) => setPort(e.target.value)}
                placeholder="端口"
              />
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="密碼"
              />
              <button onClick={handleConnect}>連接</button>
            </div>
          ) : (
            <div className="connection-info">
              <p>已連接到: {serverIP}:{port}</p>
              <p>上傳: {formatBytes(stats.bytesSent)}</p>
              <p>下載: {formatBytes(stats.bytesReceived)}</p>
              <p>連接時間: {formatDuration(stats.uptime)}</p>
              <button onClick={handleDisconnect}>斷開連接</button>
            </div>
          )}
        </div>
      )}

      {error && <div className="error">{error}</div>}
    </div>
  );
}

// Helper function to format bytes
function formatBytes(bytes) {
  if (!bytes) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`;
}

// Helper function to format duration
function formatDuration(seconds) {
  if (!seconds) return '0秒';
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const remainingSeconds = seconds % 60;
  return `${hours}時${minutes}分${remainingSeconds}秒`;
}

export default App;
