import React, { useState } from 'react';
import axios from 'axios';
import Editor from '@monaco-editor/react';

function SubmitScript() {
    const [script, setScript] = useState({
        name: '',
        code: '',
        schedule: '',
        description: '',
        alias: ''
    });

    const handleChange = (e) => {
        const { name, value } = e.target;
        setScript({ ...script, [name]: value });
    };

    const handleEditorChange = (value) => {
        setScript({ ...script, code: value });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const response = await axios.post('/api/scripts', script);
            alert(response.data.message);
        } catch (error) {
            console.error('There was an error submitting the script!', error);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="max-w-6xl mx-auto p-4 bg-white shadow-md rounded">
            <h2 className="text-xl font-bold mb-4">Submit Script</h2>
            <input type="text" name="name" placeholder="Name" onChange={handleChange} className="w-full p-2 mb-4 border rounded" />
            <Editor
                height="200px"
                defaultLanguage="lua"
                defaultValue="// Enter your Lua code here"
                value={script.code}
                onChange={handleEditorChange}
                className="w-full p-2 mb-4 border rounded"
            />
            <input type="text" name="schedule" placeholder="Schedule" onChange={handleChange} className="w-full p-2 mb-4 border rounded" />
            <input type="text" name="description" placeholder="Description" onChange={handleChange} className="w-full p-2 mb-4 border rounded" />
            <input type="text" name="alias" placeholder="Alias" onChange={handleChange} className="w-full p-2 mb-4 border rounded" />
            <button type="submit" className="w-full p-2 bg-blue-500 text-white rounded">Submit Script</button>
        </form>
    );
}

export default SubmitScript;