import React, { useRef } from 'react'

function TestUseRef() {
    const inputElement = useRef<HTMLInputElement | null>(null);

    const useFocus = () => {
        inputElement.current?.focus();
    }
  return (
    <>
        <input type="text" ref={inputElement} className="border border-gray-400 rounded px-2 py-1 focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <button onClick={useFocus} className='pl-2'>Focus</button>
    </>
  )
}

export default TestUseRef