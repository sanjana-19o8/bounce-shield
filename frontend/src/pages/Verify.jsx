import React, { useState } from "react";
import Sidebar from "../components/Sidebar";


const Verify = () => {
  const [file, setFile] = useState(null);
  const [result, setResult] = useState(null);
  const [loading, setLoading] = useState(false);

  const handleUpload = async () => {
    if (!file) return;

    const formData = new FormData();
    formData.append("emails", file);

    try {
      setLoading(true);
      const res = await fetch("/api/batch-verify", {
        method: "POST",
        body: formData,
      });
      const data = await res.json();
      setResult(data);
    } catch (err) {
      console.error("Upload failed", err);
    } finally {
      setLoading(false);
    }
  };

  const handleDownload = () => {
    if (!result || !result.length) return;

    const headers = Object.keys(result[0]);
    const csvRows = [
      headers.join(","),
      ...result.map(row => headers.map(header => row[header]).join(","))
    ];

    const blob = new Blob([csvRows.join("\n")], { type: "text/csv" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "verified_results.csv";
    a.click();
    URL.revokeObjectURL(url);
  };

  return (
    <div className='flex'>
      <Sidebar />
      <main className="p-10">
        <h2 className="text-2xl font-bold mb-4">📤 Batch Email Verification</h2>

        <input type="file" accept=".csv" onChange={(e) => setFile(e.target.files[0])} />
        <button
          className="ml-4 bg-blue-600 text-white px-4 py-2 rounded disabled:opacity-50"
          onClick={handleUpload}
          disabled={loading}
        >
          {loading ? "Verifying..." : "Upload"}
        </button>

        {result && (
          <>
            <pre className="mt-6 bg-gray-100 p-4 rounded text-sm max-h-96 overflow-auto">
              {JSON.stringify(result, null, 2)}
            </pre>
            <button
              className="mt-4 bg-green-600 text-white px-4 py-2 rounded"
              onClick={handleDownload}
            >
              ⬇ Download CSV
            </button>
          </>
        )}
      </main>
    </div>
  );
};

export default Verify;
