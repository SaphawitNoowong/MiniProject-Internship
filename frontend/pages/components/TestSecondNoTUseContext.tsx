import React from 'react'
import TestThirdNotUseContext from './TestThirdNotUseContext'

type data = {
    testData: string // หรือกำหนด type ที่เหมาะสม เช่น string, number, object, etc.
}

function TestSecondNoTUseContext({testData}: data) {
  return (
    <div>TestSecondNoTUseContext
        <hr />
        <TestThirdNotUseContext testData = {testData}/>
    </div>
  )
}

export default TestSecondNoTUseContext