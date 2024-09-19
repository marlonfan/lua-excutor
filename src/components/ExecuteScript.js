import React, { useState } from 'react';
import axios from 'axios';

function ExecuteScript() {
    const [scriptName, setScriptName] = useState('');
    const [scriptParams, setScriptParams] = useState('');
    const [result, setResult] = useState('');

    const handleExecute = async () => {
        try {
            const response = await axios.post(`/api/scripts/${scriptName}/execute`, {
                params: {
                    scriptName: scriptName,
                    scriptParams: JSON.parse(scriptParams),
                },
            });
            setResult(response.data.result);
        } catch (error) {
            console.error('Error executing script:', error);
            setResult('Error executing script');
        }
    };

    return (
        // <div className="max-w-6xl mx-auto p-4 bg-white shadow-md rounded">
            <div className="max-w-6xl mx-auto p-4 bg-white rounded-lg shadow-md">
                <h1 className="text-2xl font-bold mb-4">Execute Script</h1>
                <div className="mb-4">
                    <input
                        type="text"
                        value={scriptName}
                        onChange={(e) => setScriptName(e.target.value)}
                        placeholder="Enter script name"
                        className="w-full px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                </div>
                <div className="mb-4">
                    <textarea
                        type="text"
                        value={scriptParams}
                        onChange={(e) => setScriptParams(e.target.value)}
                        placeholder="Enter script params"
                        className="w-full px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                </div>
                <button
                    onClick={handleExecute}
                    className="w-full bg-blue-500 text-white py-2 rounded-md hover:bg-blue-600 transition duration-200"
                >
                    Execute
                </button>
                <div className="mt-6">
                    <h2 className="text-xl font-semibold mb-2">Result:</h2>
                    <pre className="bg-gray-100 p-4 rounded-md">{result}</pre>
                </div>
            </div>
        // </div>
    );
}

export default ExecuteScript;