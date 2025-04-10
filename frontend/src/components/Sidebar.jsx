import { Link } from "react-router-dom";

const Sidebar = () => {
  return (
    <div className="w-64 h-screen bg-black-200 text-white flex flex-col p-4">
      <h2 className="text-xl font-bold mb-6">🛡 BounceShield</h2>
      <Link to="/" className="mb-4 hover:text-blue-300">🏠 Dashboard</Link>
      <Link to="/verify" className="mb-4 hover:text-blue-300">📤 Batch Verify</Link>
      <Link to="/history" className="mb-4 hover:text-blue-300">🕓 Job History</Link>
    </div>
  );
};

export default Sidebar;
