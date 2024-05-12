"use client";
import { useState, useEffect, FormEvent } from 'react'

export default function Page() {
    const [userID, setUserID] = useState('')
    const [isSubmit, setIsSubmit] = useState(false)

    function onSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()
        setIsSubmit(true)
    }

    return (
     <div>
        <h1>Hello, SSE</h1>
        {!isSubmit ? (
            <div>
            <h1>Please enter ID</h1>
                <form onSubmit={onSubmit}>
                    <input type="text" name="name" id="first_name" 
                        onChange={(e) => {setUserID(e.target.value)}}
                        className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                    />
                    <button type="submit" className='text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm w-full sm:w-auto px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800'>Submit</button>
                </form>
            </div>
        ): (
            <Message userID={userID}/>
        )}
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
        <h1>Got Message:</h1>
        <h1>{data}</h1>
      </div>
    )
}

