import React, { useState } from 'react';
import axios from 'axios';
import Editor from '@monaco-editor/react';

function EditScript({ script, onClose, onUpdate }) {
    const [updatedScript, setUpdatedScript] = useState(script);

    const handleChange = (e) => {
        const { name, value } = e.target;
        setUpdatedScript({ ...updatedScript, [name]: value });
    };

    const handleEditorChange = (value) => {
        setUpdatedScript({ ...updatedScript, code: value });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const response = await axios.put(`/api/update/${script.name}`, updatedScript);
            alert(response.data.message);
            onUpdate();
            onClose();
        } catch (error) {
            console.error('There was an error updating the script!', error);
        }
    };

    return (
        <div className="fixed inset-0 bg-gray-600 bg-opacity-50 flex justify-center items-center">
            <div className="bg-white p-4 rounded shadow-md w-3/4">
                <h2 className="text-2xl font-bold mb-4">Edit Script</h2>
                <form onSubmit={handleSubmit}>
                    <input type="text" name="name" value={updatedScript.name} onChange={handleChange} className="w-full p-2 mb-4 border rounded" readOnly />
                    <Editor
                        height="200px"
                        defaultLanguage="lua"
                        value={updatedScript.code}
                        onChange={handleEditorChange}
                        className="w-full p-2 mb-4 border rounded"
                    />
                    <input type="text" name="schedule" value={updatedScript.schedule} onChange={handleChange} className="w-full p-2 mb-4 border rounded" />
                    <input type="text" name="description" value={updatedScript.description} onChange={handleChange} className="w-full p-2 mb-4 border rounded" />
                    <input type="text" name="alias" value={updatedScript.alias} onChange={handleChange} className="w-full p-2 mb-4 border rounded" />
                    <div className="flex justify-end space-x-4">
                        <button type="button" onClick={onClose} className="p-2 bg-gray-500 text-white rounded">Cancel</button>
                        <button type="submit" className="p-2 bg-blue-500 text-white rounded">Update Script</button>
                    </div>
                </form>
            </div>
        </div>
    );
}

export default EditScript;