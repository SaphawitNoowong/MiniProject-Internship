import React, { useMemo, useState } from 'react'

function TryUseMeMo() {
    const [count, setCount] = useState(0)
    const [text, setText] = useState("")

    const expensiveValue = useMemo(() => {
        console.log("New calcurate")
        let total = 0
        for (let i = 0; i < 100; i++) {
            total += 1
        }
        return count + total
    }, [count])
    return (
        <div>
            <h1>Try UseMemo</h1>
            <div>Expensive value: {expensiveValue}</div>
            <button onClick={() => setCount(count + 1)}
                className="border border-gray-400 rounded px-2 py-1 focus:outline-none focus:ring-2 focus:ring-blue-500 m-2">
                Count
            </button>
            <input type="text" value={text} onChange={(e) => setText(e.target.value)} placeholder='Try input text'
                className="border border-gray-400 rounded px-2 py-1 focus:outline-none focus:ring-2 focus:ring-blue-500"/>
            <div>
                Text : {text}
            </div>
        </div>
    )
}

export default TryUseMeMo