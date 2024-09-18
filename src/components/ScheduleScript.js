import React, { useState } from 'react';
import axios from 'axios';

function ScheduleScript() {
    const [script, setScript] = useState({
        name: '',
        schedule: ''
    });

    const handleChange = (e) => {
        const { name, value } = e.target;
        setScript({ ...script, [name]: value });
    };

    const handleSchedule = async (e) => {
        e.preventDefault();
        try {
            const response = await axios.post(`/api/scripts/${script.name}/schedule`, script);
            alert(response.data.message);
        } catch (error) {
            console.error('There was an error scheduling the script!', error);
        }
    };

    return (
        <form onSubmit={handleSchedule} className="max-w-6xl mx-auto p-4 bg-white shadow-md rounded">
            <h2 className="text-xl font-bold mb-4">Schedule Script</h2>
            <input type="text" name="name" placeholder="Name" onChange={handleChange} className="w-full p-2 mb-4 border rounded" />
            <input type="text" name="schedule" placeholder="Schedule" onChange={handleChange} className="w-full p-2 mb-4 border rounded" />
            <button type="submit" className="w-full p-2 bg-blue-500 text-white rounded">Schedule Script</button>
        </form>
    );
}

export default ScheduleScript;