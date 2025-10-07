import React, {useContext} from 'react'
import { dataContext } from '../index'
function ThirdComponent() {
    const data = useContext(dataContext);
  return (
    <div>ThirdComponent
        <hr />
        <div>From Context : {data}</div>
    </div>
  )
}

export default ThirdComponent