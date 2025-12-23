import { useParams, useNavigate } from "react-router-dom";

export default function Topbar() {
  const { orgId } = useParams();
  const navigate = useNavigate();

  function logout() {
    fetch("/api/logout", { method: "POST" }).finally(() => {
      navigate("/login");
    });
  }

  return (
    <header className="h-14 border-b flex items-center justify-between px-6 bg-white">
      {/* Left */}
      <div className="flex items-center gap-4">
        <span className="font-semibold text-lg">Valora</span>
        <span className="text-sm text-gray-500">Org: {orgId}</span>
      </div>

      {/* Right */}
      <div className="flex items-center gap-4">
        <button
          className="text-sm text-gray-600 hover:text-black"
          onClick={() => navigate("/orgs")}
        >
          Switch Org
        </button>

        <div className="relative">
          <button className="text-sm font-medium">Account ▾</button>
          {/* simple dropdown placeholder */}
        </div>

        <button
          onClick={logout}
          className="text-sm text-red-600 hover:text-red-800"
        >
          Logout
        </button>
      </div>
    </header>
  );
}
