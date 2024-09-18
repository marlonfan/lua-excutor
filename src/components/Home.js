import React, { useState, useEffect } from 'react';
import axios from 'axios';
import EditScript from './EditScript';

function Home() {
    const [scripts, setScripts] = useState([]);
    const [editingScript, setEditingScript] = useState(null);

    useEffect(() => {
        const fetchScripts = async () => {
            try {
                const response = await axios.get('/api/scripts');
                setScripts(response.data);
            } catch (error) {
                console.error('There was an error fetching the scripts!', error);
            }
        };

        fetchScripts();
    }, []);

    const handleEdit = (script) => {
        setEditingScript(script);
    };

    const handleUpdate = () => {
        const fetchScripts = async () => {
            try {
                const response = await axios.get('/api/scripts');
                setScripts(response.data);
            } catch (error) {
                console.error('There was an error fetching the scripts!', error);
            }
        };

        fetchScripts();
    };

    return (
        <div className="max-w-6xl mx-auto p-4 bg-white shadow-md rounded">
            <h2 className="text-xl font-bold mb-4">Registered Scripts</h2>
            <table className="min-w-full bg-white">
                <thead>
                    <tr>
                        <th className="py-2 px-4 text-left">Name</th>
                        <th className="py-2 px-4 text-left">Description</th>
                        <th className="py-2 px-4 text-left">Alias</th>
                        <th className="py-2 px-4 text-left">Schedule</th>
                        <th className="py-2 px-4 text-left">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {scripts.map(script => (
                        <tr key={script.name} className="border-t">
                            <td className="py-2 px-4">{script.name}</td>
                            <td className="py-2 px-4">{script.description}</td>
                            <td className="py-2 px-4">{script.alias}</td>
                            <td className="py-2 px-4">{script.schedule}</td>
                            <td className="py-2 px-4">
                                <button onClick={() => handleEdit(script)} className="p-2 bg-blue-500 text-white rounded">Edit</button>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
            {editingScript && <EditScript script={editingScript} onClose={() => setEditingScript(null)} onUpdate={handleUpdate} />}
        </div>
    );
}

export default Home;