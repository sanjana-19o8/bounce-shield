import { useState } from "react";
import axios from "axios";
import toast from "react-hot-toast";
import Sidebar from "../components/Sidebar";
import { CheckCircle, MailCheck } from "lucide-react";

const Dashboard = ({ user }) => {
  const [email, setEmail] = useState("");
  const [result, setResult] = useState(null);

  const handleVerify = async () => {
    if (!email) return toast.error("Please enter an email!");
    try {
      const res = await axios.post("/api/verify", { email });
      setResult(res.data);
      toast.success("✅ Email verification complete!");
    } catch {
      toast.error("Something went wrong while verifying.");
    }
  };

  return (
    <div className="flex min-h-screen bg-gray-50">
      <Sidebar />
      <main className="flex-1 p-6 sm:p-10">
        <h2 className="text-3xl font-bold mb-6 text-gray-800">📬 Dashboard</h2>

        <div className="max-w-xl w-full bg-white shadow-2xl rounded-2xl p-6 sm:p-8 mx-auto">
          <div className="flex items-center gap-2 mb-4">
            <MailCheck className="text-blue-600" size={24} />
            <h3 className="text-xl font-semibold text-gray-700">
              Check New Email
            </h3>
          </div>

          <div className="flex flex-col sm:flex-row sm:items-center gap-3">
            <input
              type="email"
              placeholder="example@domain.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="flex-1 border border-gray-300 px-4 py-2 rounded-lg text-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-400"
            />
            <button
              onClick={handleVerify}
              className="bg-blue-600 hover:bg-blue-700 text-white font-medium px-5 py-2 rounded-lg transition-all shadow"
            >
              Verify
            </button>
          </div>

          {result && (
            <div className="mt-6 bg-gray-100 border border-gray-200 rounded-xl p-4 space-y-2">
              <div className="flex items-center gap-2">
                <CheckCircle className="text-green-600" size={20} />
                <p className="text-sm text-gray-700">
                  <strong>Email:</strong> {result.email}
                </p>
              </div>
              <p className="text-sm text-gray-700">
                <strong>Status:</strong>{" "}
                <span
                  className={`${
                    result.status === "valid"
                      ? "text-green-600"
                      : "text-red-500"
                  } font-semibold`}
                >
                  {result.status}
                </span>
              </p>
              <p className="text-sm text-gray-700">
                <strong>Reason:</strong> {result.reason}
              </p>
            </div>
          )}
        </div>
      </main>
    </div>
  );
};

export default Dashboard;
