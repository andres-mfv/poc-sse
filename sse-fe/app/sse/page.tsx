"use client";
import { useState, useEffect } from 'react'

export default function Page() {
    return (
     <div>
        <h1>Hello, SSE</h1>
        <Message />
    </div>
    )
}

function Message() {
    const [data, setData] = useState('')
    const [isLoading, setLoading] = useState(true)
   
    useEffect(() => {
        const source = new EventSource(`http://localhost:8080/sse`);

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
   
    if (isLoading) return <p>Loading...</p>
    if (!data) return <p>No profile data</p>
   
    return (
      <div>
        <h1>{data}</h1>
        <p>{data}</p>
      </div>
    )
  }
