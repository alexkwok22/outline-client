import { useState, useEffect } from 'react';
import './App.css';
import { Connect, Disconnect, GetStatus, GetStats, ActivateLicense, GetLicenseStatus } from '../wailsjs/go/main/App';

function App() {
  const [serverIP, setServerIP] = useState('');
  const [port, setPort] = useState('');
  const [password, setPassword] = useState('');
  const [isConnected, setIsConnected] = useState(false);
  const [stats, setStats] = useState({});
  const [licenseKey, setLicenseKey] = useState('');
  const [isLicenseValid, setIsLicenseValid] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    // Check initial connection status
    updateStatus();
    // Check license status
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
      checkLicense();
    } catch (err) {
      setError(err.toString());
    }
  };

  return (
    <div className="container">
      <h1>Outline VPN Client</h1>
      
      {!isLicenseValid && (
        <div className="license-section">
          <h2>Activate License</h2>
          <input
            type="text"
            placeholder="Enter License Key"
            value={licenseKey}
            onChange={(e) => setLicenseKey(e.target.value)}
          />
          <button onClick={handleActivateLicense}>Activate</button>
        </div>
      )}

      {isLicenseValid && (
        <div className="vpn-section">
          <h2>VPN Connection</h2>
          <div className="input-group">
            <input
              type="text"
              placeholder="Server IP"
              value={serverIP}
              onChange={(e) => setServerIP(e.target.value)}
              disabled={isConnected}
            />
            <input
              type="text"
              placeholder="Port"
              value={port}
              onChange={(e) => setPort(e.target.value)}
              disabled={isConnected}
            />
            <input
              type="password"
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              disabled={isConnected}
            />
          </div>

          <div className="button-group">
            {!isConnected ? (
              <button onClick={handleConnect}>Connect</button>
            ) : (
              <button onClick={handleDisconnect}>Disconnect</button>
            )}
          </div>

          {isConnected && (
            <div className="stats">
              <h3>Connection Statistics</h3>
              <p>Bytes Received: {stats.bytesReceived || 0}</p>
              <p>Bytes Sent: {stats.bytesSent || 0}</p>
              <p>Uptime: {stats.uptime || 0}s</p>
            </div>
          )}
        </div>
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

export default App;
