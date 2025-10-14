import { useQuery } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import Header from "./components/Header"; 
import CreateButtonNisit from "./components/CreateButtonNisit";

export default function Home() {
  const { data, isLoading, isError } = useQuery({
    queryKey: ["users"],
    queryFn: async () => {
      const res = await fetch("http://localhost:5000/users");
      if (!res.ok) throw new Error("Failed to fetch users");
      return res.json();
    },
  });

  const [showDetails, setShowDetails] = useState(false);

  useEffect(() => {
    if (data?.data) {
      console.log(`Have ${data.data.length} user in database`);
    }
  }, [showDetails, data]); // ใส่ data ใน dependency array ถูกต้องแล้วครับ

  return (
    // ให้ div หลักเป็นแค่พื้นหลัง ไม่ต้องมี flex สำหรับ layout หลัก
    <div className="bg-gray-50 min-h-screen">
      {/* 1. วาง Header ไว้บนสุดได้เลย */}
      <Header />

      {/* 2. เนื้อหาหลักจะอยู่ข้างล่าง Header โดยอัตโนมัติ */}
      <main className="p-6 max-w-4xl mx-auto">
        <h1 className="text-2xl font-bold mb-6 text-gray-800">Student Dashboard</h1>

        {isLoading && <p>Loading...</p>}
        {isError && <p>Failed to load</p>}

        <div className="flex flex-col gap-4">
          <button
            type="button"
            onClick={() => setShowDetails(!showDetails)}
            className="self-start px-4 py-2 rounded bg-blue-600 text-white hover:bg-blue-700 transition-colors"
          >
            {showDetails ? "Hide Users" : "Show All Users"}
          </button>
          <CreateButtonNisit />

          {showDetails && !isLoading && !isError && (
            <ul className="space-y-3">
              {(data?.data ?? []).map((u: any, i: number) => (
                <li key={i} className="rounded border bg-white p-4 shadow-sm hover:shadow-md transition-shadow">
                  {/* ... User details ... */}
                  <div className="text-sm text-gray-500">
                    StudentCode :{" "}
                    <span className="font-medium text-gray-900">{u.studentCode}</span>
                  </div>
                  <div className="text-sm text-gray-500 mt-2">
                    Name :{" "}
                    <span className="font-medium text-gray-900">{u.name}</span>
                  </div>
                  <div className="text-sm text-gray-500 mt-2">
                    Major :{" "}
                    <span className="font-medium text-gray-900">{u.major}</span>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </div>
      </main>
    </div>
  );
}