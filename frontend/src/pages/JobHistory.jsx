import { useEffect, useState } from "react";
import axios from "axios";
import toast from "react-hot-toast";
import Sidebar from "../components/Sidebar";

const JobHistory = ({ user }) => {
  const [jobs, setJobs] = useState([]);

  useEffect(() => {
    axios
      .get(`/api/jobs?user=${user}`)
      .then((res) => setJobs(res.data))
      .catch(() => toast.error("Failed to fetch job history"));
  }, [user]);

  return (
    <div className="flex">
      <Sidebar />
      <main className="p-6 w-full">
        <h2 className="text-2xl font-bold mb-4">🕓 Job History</h2>
        <div className="space-y-6">
          {jobs.length === 0 && <p>No jobs yet.</p>}
          {jobs.map((job) => (
            <div key={job.id} className="border rounded-xl p-4 bg-white shadow">
              <p className="font-semibold text-sm text-gray-600">
                Job ID: <span className="text-gray-800">{job.id}</span>
              </p>
              <p className="text-sm text-gray-600">
                Timestamp: {new Date(job.timestamp * 1000).toLocaleString()}
              </p>
              <ul className="mt-2 space-y-1 text-sm text-gray-700">
                {job.results.map((r, idx) => (
                  <li key={idx}>
                    <span className="font-medium">{r.email}</span> -{" "}
                    <span className={`text-${r.status === "valid" ? "green" : "red"}-500`}>
                      {r.status}
                    </span>{" "}
                    ({r.reason})
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>
      </main>
    </div>
  );
};

export default JobHistory;
