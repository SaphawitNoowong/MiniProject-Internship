import React, { useCallback, useState } from 'react'
import TryUseCallback from './TryUseCallback'
function TryUsecallbackParent() {
    const [count, setCount] = useState(0)
    const [text, setText] = useState("")

    const handleClick = useCallback(() => {
        setCount((c) => c + 1)
    }, [])

    return (
        <div>
            <h1>Try useCallback</h1>
            <div>Count : {count}</div>
            <div>
                <TryUseCallback onClick={handleClick} />
            </div>
            <input type="text" value={text} onChange={(e) => setText(e.target.value)} placeholder='Try input Text'
                className="border border-gray-400 rounded px-2 py-1 focus:outline-none focus:ring-2 focus:ring-blue-500" />
            <div>Text : {text}</div>
        </div>
    )
}

export default TryUsecallbackParent