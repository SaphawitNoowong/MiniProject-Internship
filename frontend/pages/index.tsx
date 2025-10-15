import { useQuery } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import Header from "./components/Header";
import CreateButtonNisit from "./components/CreateButtonNisit";
import UpdateButtonNisit from "./components/UpdateButtonNisit";
import DeleteButtonNisit from "./components/DeleteButtonNisit";

export default function Home() {
  const { data, isLoading, isError } = useQuery({
    queryKey: ["users"],
    queryFn: async () => {
      const res = await fetch("http://localhost:5000/users");
      if (!res.ok) throw new Error("Failed to fetch users");
      return res.json();
    },
  });


  return (
    <div className="bg-gray-50 min-h-screen">
      <Header />
      <main className="p-6 max-w-4xl mx-auto">

        {isLoading && <p>Loading...</p>}
        {isError && <p>Failed to load</p>}

        <div className="flex flex-col gap-4">
          <div className="flex justify-end">
            <CreateButtonNisit />
          </div>
          {!isLoading && !isError && (
            <ul className="space-y-3">
              {(data?.data ?? []).map((u: any) => (
                <li key={u.studentCode} className="rounded border bg-white p-4 shadow-sm hover:shadow-md transition-shadow">
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
                  <div className="flex justify-start md:justify-end mt-3 md:mt-0">
                    <div className="mr-5">
                    <DeleteButtonNisit initialData={u} />
                    </div>
                    <div>
                    <UpdateButtonNisit initialData={u} />
                    </div>
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