import React, { useState } from 'react';

const IPInfoApp = () => {
  const [ipAddress, setIpAddress] = useState('');
  const [ipInfo, setIpInfo] = useState(null);
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setIpInfo(null);

    try {
      const response = await fetch(`http://localhost:8080/api/ipinfo?ip=${ipAddress}`);
      if (!response.ok) {
        throw new Error('Failed to fetch IP information');
      }
      const data = await response.json();
      setIpInfo(data);
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">IP Address Information</h1>
      <form onSubmit={handleSubmit} className="mb-4">
        <input
          type="text"
          value={ipAddress}
          onChange={(e) => setIpAddress(e.target.value)}
          placeholder="Enter IP address"
          className="border p-2 mr-2"
        />
        <button type="submit" className="bg-blue-500 text-white p-2 rounded">
          Get Info
        </button>
      </form>
      {error && <p className="text-red-500">{error}</p>}
      {ipInfo && (
        <div className="bg-gray-100 p-4 rounded">
          <h2 className="text-xl font-semibold mb-2">IP Information:</h2>
          <p><strong>IP Address:</strong> {ipInfo.ip}</p>
          <p><strong>Subnet:</strong> {ipInfo.subnet}</p>
          <p><strong>Gateway:</strong> {ipInfo.gateway}</p>
          <p><strong>Class:</strong> {ipInfo.class}</p>
          <p><strong>Type:</strong> {ipInfo.is_private ? 'Private' : 'Public'}</p>
        </div>
      )}
    </div>
  );
};

export default IPInfoApp;
