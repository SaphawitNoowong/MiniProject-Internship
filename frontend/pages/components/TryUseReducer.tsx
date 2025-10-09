import React, { useReducer } from 'react'
type State = {
    number: number;
};

type Action =
    | { type: 'increment' }
    | { type: 'decrement' }
    | { type: 'reset' };

function TryUseReducer() {
    const [state, dispatch] = useReducer(reducer, 0, init);
    return (
        <div>
            <h1> Number = {state.number}</h1>
            <button onClick={() => dispatch({ type: "increment" })}
                className="border border-gray-400 rounded px-2 py-1 focus:outline-none focus:ring-2 focus:ring-blue-500 p-2 m-2">
                Increment
            </button>
            <button onClick={() => dispatch({ type: "decrement" })}
                className="border border-gray-400 rounded px-2 py-1 focus:outline-none focus:ring-2 focus:ring-blue-500 p-2 m-2">
                Decrement
            </button>
            <button onClick={() => dispatch({ type: "reset" })}
                className="border border-gray-400 rounded px-2 py-1 focus:outline-none focus:ring-2 focus:ring-blue-500p-2 m-2">
                Reset
            </button>
        </div>
    )
}

function reducer(state: State, action: Action) {
    switch (action.type) {
        case "increment":
            return { number: state.number + 1 };
        case "decrement":
            return { number: state.number - 1 };
        case "reset":
            return init(0)
        default:
            return state;
    }
}

function init(initialNumber: number) {
    return { number: initialNumber };
}

export default TryUseReducer