import React from 'react'

type data = {
    testData: string // หรือกำหนด type ที่เหมาะสม เช่น string, number, object, etc.
}

function TestThirdNotUseContext({testData}: data) {
  return (
    <div>TestThirdNotUseContext
        <hr />
        <div>Not use context : {testData}</div>
    </div>
  )
}

export default TestThirdNotUseContext