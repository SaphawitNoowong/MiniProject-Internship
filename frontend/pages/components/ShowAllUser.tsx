import React from 'react'
import { useQuery, keepPreviousData } from "@tanstack/react-query";
import { useState, useEffect } from "react";
import Header from './Header';
import CreateButtonNisit from "./CreateButtonNisit";
import UpdateButtonNisit from "./UpdateButtonNisit";
import DeleteButtonNisit from "./DeleteButtonNisit";


function ShowAllUser() {
    const [page, setPage] = useState(1);
    const limit = 10;
    const [searchTerm, setSearchTerm] = useState("");
    // State สำหรับคำค้นหาที่ผ่านการ Debounce แล้ว
    const [debouncedSearchTerm, setDebouncedSearchTerm] = useState("");


    // 2. ใช้ useEffect ทำ Debouncing
    useEffect(() => {
        // ตั้งเวลา 500ms ก่อนที่จะอัปเดต debouncedSearchTerm
        const handler = setTimeout(() => {
            setDebouncedSearchTerm(searchTerm);
            setPage(1); //เมื่อมีการค้นหาใหม่ ให้กลับไปที่หน้า 1 เสมอ
        }, 500); // 500ms delay

        // ยกเลิก timer เก่า เพื่อเริ่มนับเวลาใหม่
        return () => {
            clearTimeout(handler);
        };
    }, [searchTerm]);

    const { data, isLoading, isError, isFetching } = useQuery({
        queryKey: ["users", page, debouncedSearchTerm],
        queryFn: async () => {
            const res = await fetch(`http://localhost:5000/users?page=${page}&limit=${limit}&search=${debouncedSearchTerm}`);
            if (!res.ok) throw new Error("Failed to fetch users");
            return res.json();
        },
        // (แนะนำ) ทำให้ข้อมูลเก่าแสดงค้างไว้ระหว่างรอโหลดหน้าใหม่ ป้องกันหน้ากระพริบ
        placeholderData: keepPreviousData,
    });
    return (
        <div className="bg-gray-50 min-h-screen">
            <Header>
                {/* UI ส่วนนี้จะถูกนำไปแสดงใน "ช่อง" ที่เราเตรียมไว้ใน Header */}
                <div className="flex items-center w-full">
                    <input
                        type="text"
                        placeholder="Search by student code name or major"
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        className="px-4 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                    />

                </div>
            </Header>
            <main className="p-6 max-w-4xl mx-auto">

                {/* แสดง Loading แค่ครั้งแรกที่โหลด */}
                {isLoading && <p className="text-center">Loading...</p>}
                {isError && <p className="text-center">Failed to load data.</p>}

                {!isLoading && !isError && (
                    <div className="flex flex-col gap-4">
                        <div className='flex justify-end'>
                         <CreateButtonNisit />
                         </div>
                        {/* แสดงข้อมูลจาก data.data */}
                        <ul className="space-y-3">
                            {(data?.data ?? []).map((u: any) => (
                                <li key={u.studentCode} className="rounded border bg-white p-4 shadow-sm hover:shadow-md transition-shadow">
                                    {/* ... User details ... */}
                                    <div className="flex items-center justify-between">
                                        <div>
                                            <span className="text-sm text-gray-500">StudentCode :{" "}</span>
                                            <span className="text-sm font-medium text-gray-900">{u.studentCode}</span>
                                        </div>
                                        <div>
                                            <DeleteButtonNisit initialData={u} />
                                        </div>
                                    </div>
                                    <div className="text-sm text-gray-500 mt-2">
                                        Name :{" "}
                                        <span className="font-medium text-gray-900">{u.name}</span>
                                    </div>
                                    <div className="flex items-center justify-between">
                                        <div className="text-sm text-gray-500 mt-2">
                                            Major :{" "}
                                            <span className="font-medium text-gray-900">{u.major}</span>
                                        </div>
                                        <UpdateButtonNisit initialData={u} />
                                    </div>
                                </li>
                            ))}
                        </ul>

                        {/* Pagination Controls */}
                        <div className="flex items-center justify-center space-x-4 mt-3">
                            <button
                                onClick={() => setPage(old => Math.max(old - 1, 1))}
                                disabled={page === 1}
                                className="px-4 py-2 rounded bg-white border border-gray-300 text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                            >
                                Previous
                            </button>

                            <span>
                                Page {data?.pagination.currentPage} of {data?.pagination.totalPages}
                            </span>

                            <button
                                onClick={() => setPage(old => old + 1)}
                                // Disable ปุ่ม Next ถ้าหน้าปัจจุบันคือหน้าสุดท้าย
                                disabled={page === data?.pagination.totalPages || !data?.data?.length}
                                className="px-4 py-2 rounded bg-white border border-gray-300 text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                            >
                                Next
                            </button>
                        </div>

                        {/* (Bonus) แสดงสถานะขณะกำลังโหลดหน้าถัดไป */}
                        {/* {isFetching && <div className="text-center text-gray-500">Fetching next page...</div>} */}
                    </div>
                )}
            </main>
        </div>
    )
}

export default ShowAllUser