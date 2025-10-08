import React, { useEffect, useRef, useState } from 'react'

function TestUseRefConst() {
    const count = useRef(0);
    const [inputValue, setInputValue] = useState("");
    console.log(count.current)
    
    useEffect(() => {
        count.current = count.current + 1 
    },[inputValue])
  return (
    <div>
        <input type="text" value={inputValue} onChange={(e) => setInputValue(e.target.value)}
          className="border border-gray-400 rounded px-2 py-1 focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <h1>Count Word : {count.current}</h1>
    </div>
  )
}

export default TestUseRefConst