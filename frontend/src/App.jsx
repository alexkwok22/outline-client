import { useState, useEffect } from 'react';
import './App.css';
import { Connect, Disconnect, GetStatus, GetStats, ActivateLicense, GetLicenseStatus, GetLicenseInfo } from '../wailsjs/go/main/App';

function App() {
  const [serverIP, setServerIP] = useState('');
  const [port, setPort] = useState('');
  const [password, setPassword] = useState('');
  const [isConnected, setIsConnected] = useState(false);
  const [stats, setStats] = useState({});
  const [licenseKey, setLicenseKey] = useState('');
  const [isLicenseValid, setIsLicenseValid] = useState(false);
  const [licenseInfo, setLicenseInfo] = useState({});
  const [error, setError] = useState('');

  useEffect(() => {
    // Check initial connection status
    updateStatus();
    // Check license status and info
    checkLicense();
    // Update stats periodically
    const interval = setInterval(updateStats, 1000);
    return () => clearInterval(interval);
  }, []);

  const updateStatus = async () => {
    try {
      const status = await GetStatus();
      setIsConnected(status);
    } catch (err) {
      setError(err.toString());
    }
  };

  const updateStats = async () => {
    if (isConnected) {
      try {
        const stats = await GetStats();
        setStats(stats);
      } catch (err) {
        setError(err.toString());
      }
    }
  };

  const checkLicense = async () => {
    try {
      const status = await GetLicenseStatus();
      setIsLicenseValid(status);
      if (status) {
        const info = await GetLicenseInfo();
        setLicenseInfo(info);
      }
    } catch (err) {
      setError(err.toString());
    }
  };

  const handleConnect = async () => {
    try {
      await Connect(serverIP, parseInt(port), password);
      updateStatus();
    } catch (err) {
      setError(err.toString());
    }
  };

  const handleDisconnect = async () => {
    try {
      await Disconnect();
      updateStatus();
    } catch (err) {
      setError(err.toString());
    }
  };

  const handleActivateLicense = async () => {
    try {
      await ActivateLicense(licenseKey);
      await checkLicense();
      setLicenseKey(''); // Clear input after successful activation
    } catch (err) {
      setError(err.toString());
    }
  };

  return (
    <div className="container">
      <h1>Outline VPN Client</h1>
      
      {!isLicenseValid ? (
        <div className="license-section">
          <h2>授權啟動</h2>
          <div className="input-group">
            <input
              type="text"
              placeholder="請輸入授權碼"
              value={licenseKey}
              onChange={(e) => setLicenseKey(e.target.value)}
            />
            <button onClick={handleActivateLicense}>啟動</button>
          </div>
        </div>
      ) : (
        <>
          <div className="license-info">
            <h3>授權資訊</h3>
            <p>狀態: {licenseInfo.status}</p>
            <p>使用者: {licenseInfo.issuedTo}</p>
            <p>到期時間: {licenseInfo.expiresIn}</p>
          </div>

          <div className="vpn-section">
            <h2>VPN 連接</h2>
            <div className="input-group">
              <input
                type="text"
                placeholder="伺服器 IP"
                value={serverIP}
                onChange={(e) => setServerIP(e.target.value)}
                disabled={isConnected}
              />
              <input
                type="text"
                placeholder="端口"
                value={port}
                onChange={(e) => setPort(e.target.value)}
                disabled={isConnected}
              />
              <input
                type="password"
                placeholder="密碼"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                disabled={isConnected}
              />
            </div>

            <div className="button-group">
              {!isConnected ? (
                <button onClick={handleConnect}>連接</button>
              ) : (
                <button onClick={handleDisconnect}>斷開連接</button>
              )}
            </div>

            {isConnected && (
              <div className="stats">
                <h3>連接統計</h3>
                <p>已接收: {formatBytes(stats.bytesReceived || 0)}</p>
                <p>已發送: {formatBytes(stats.bytesSent || 0)}</p>
                <p>連接時長: {formatDuration(stats.uptime || 0)}</p>
              </div>
            )}
          </div>
        </>
      )}

      {error && (
        <div className="error">
          {error}
          <button onClick={() => setError('')}>✕</button>
        </div>
      )}
    </div>
  );
}

// Helper function to format bytes
function formatBytes(bytes) {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${(bytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`;
}

// Helper function to format duration
function formatDuration(seconds) {
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const remainingSeconds = seconds % 60;
  return `${hours}:${minutes.toString().padStart(2, '0')}:${remainingSeconds.toString().padStart(2, '0')}`;
}

export default App;
