"use client";
import { useState, useEffect } from 'react'

export default function Page() {
    const [userID, setUserID] = useState('')
    return (
     <div>
        <h1>Hello, SSE</h1>
        <h1>Got Message:</h1>
        <Message userID="19"/>
    </div>
    )
}

function Message(props) {
    const [data, setData] = useState('')
    const [isLoading, setLoading] = useState(true)
    console.log('userID', props.userID)
   
    useEffect(() => {
        const source = new EventSource(`http://localhost:8080/sse?user_id=${props.userID}`);

        source.addEventListener('open', () => {
            console.log('SSE opened!');
            setLoading(false);
        });
    
        source.addEventListener('message', (e) => {
            console.log(e.data);
            setData(e.data);
        });
    
        source.addEventListener('error', (e) => {
            console.error('Error: ',  e);
        });
    }, [])
   
    return (
      <div>
        <h1>{data}</h1>
      </div>
    )
}

