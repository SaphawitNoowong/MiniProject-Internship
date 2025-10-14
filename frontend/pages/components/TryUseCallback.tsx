import React from 'react'

type ChildProps = {
    onClick: () => void;
};

const TryUseCallback = React.memo(({ onClick }: ChildProps) => {
    console.log("Child Render")
    return (
        <div>
            <button onClick={onClick}
                className="border border-gray-400 rounded px-2 py-1 focus:outline-none focus:ring-2 focus:ring-blue-500">
                Count
            </button>
        </div>
    )
})

export default TryUseCallback