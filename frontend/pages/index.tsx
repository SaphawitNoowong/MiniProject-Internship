import { useQuery } from "@tanstack/react-query";
import { useEffect, useState, createContext  } from "react";

import FirstComponent from "./components/FirstComponent";
import TestNotUseContext from "./components/TestNotUseContext";
export const dataContext = createContext<string>("");

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
  const [testData, setTestData] = useState("Hi Everyone")

  let num 
  useEffect(() => {
    if (data?.data) {
      console.log(`Have ${data.data.length} user in database`);
      num = 10
      console.log("Num is ",num)
    }
    return () => {
      num = 0
      console.log("Num is ",num)
    };
  }, [showDetails]);


  return (
    <div className="p-6 max-w-3xl mx-auto">
      <h1 className="text-2xl font-bold mb-4">Users from MongoDB</h1>

      {isLoading && <p>Loading...</p>}
      {isError && <p>Failed to load</p>}
      <div className="flex flex-col gap-4">
        <button
          type="button"
          onClick={() => setShowDetails(!showDetails)}
          className="self-start px-4 py-2 rounded bg-blue-600 text-white hover:bg-blue-700"
        >
          {showDetails ? "Hide user" : "Show all user"}
        </button>

        {showDetails && !isLoading && !isError && (
          <ul className="space-y-3">
            {(data?.data ?? []).map((u: any, i: number) => (
              <li key={i} className="rounded border p-3">
                <div className="text-sm text-gray-500">StudentCode : <span> </span>
                  <span className="font-medium">{u.studentCode}</span>
                </div>
                <div className="text-sm text-gray-500 mt-2">Name : <span> </span> 
                  <span className="font-medium">{u.name}</span>
                </div>
                <div className="text-sm text-gray-500 mt-2">Major : <span> </span>
                  <span className="font-medium">{u.major}</span>
                </div>
              </li>
            ))}
          </ul>
        )}
        <dataContext.Provider value={testData}>
          <FirstComponent />
        </dataContext.Provider>
        <div>
          <TestNotUseContext testData={testData}/>
        </div>
      </div>
    </div>
  );
}
