import React from 'react';

export default function JobItem({ job }) {
  return (
    <li className="border-b py-2 text-sm">
      <strong>{job.filename}</strong> – {job.status} <br />
      <span className="text-xs text-gray-500">{job.timestamp}</span>
    </li>
  );
}
