import React from 'react'
import TestSecondNoTUseContext from './TestSecondNoTUseContext'

type data = {
    testData: string // หรือกำหนด type ที่เหมาะสม เช่น string, number, object, etc.
}

function TestNotUseContext({testData}: data) {
  return (
    <div>TestNotUseContext
        <hr />
        <TestSecondNoTUseContext  testData= {testData}/>
    </div>

  )
}

export default TestNotUseContext