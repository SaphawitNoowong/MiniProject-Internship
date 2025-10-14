import React, { useEffect, useRef, useState } from 'react'

function TestUseRefConst() {
    const [inputValue, setInputValue] = useState("");

    // count.current จะถูกเก็บค่าไว้ตลอด life cycle ของ component
    const changeCount = useRef(0);

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        // 1. อัปเดต state เพื่อให้ input แสดงผลค่าใหม่ (ทำให้เกิด re-render)
        setInputValue(e.target.value);

        // 2. อัปเดตค่าใน ref โดยตรง
        // การทำแบบนี้จะไม่ทำให้ component re-render อีกรอบ
        changeCount.current = changeCount.current + 1;

        console.log("Input changed! New count:", changeCount.current);
    };

    return (
        <div>
            <input  
                type="text"
                value={inputValue}
                onChange={handleInputChange} // เรียกใช้ฟังก์ชัน handleInputChange
                className="border border-gray-400 rounded px-2 py-1 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            {/* ข้อควรระวัง: ค่าที่แสดงผลที่นี่จะไม่อัปเดตทันที
        เพราะการเปลี่ยน .current ไม่ทำให้ re-render
        ค่าจะอัปเดตบนหน้าจอเมื่อมี re-render จากสาเหตุอื่น (เช่นการพิมพ์ใน input)
      */}
            <h1>Change Count: {changeCount.current}</h1>
        </div>
    );
}

export default TestUseRefConst