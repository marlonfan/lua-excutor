import React from 'react';
import { BrowserRouter as Router, Route, Switch, Link } from 'react-router-dom';
import Home from './components/Home';
import SubmitScript from './components/SubmitScript';
import ExecuteScript from './components/ExecuteScript';
import ScheduleScript from './components/ScheduleScript';

function App() {
    return (
        <Router>
            <div className="min-h-screen bg-gray-100">
                <nav className="bg-white shadow-md fixed w-full z-10">
                    <div className="container mx-auto px-4 py-2 flex justify-between items-center">
                        <div className="text-xl font-bold text-blue-500">Script Manager</div>
                        <ul className="flex space-x-4">
                            <li>
                                <Link to="/" className="text-gray-700 hover:text-blue-500 font-semibold">
                                    Home
                                </Link>
                            </li>
                            <li>
                                <Link to="/submit" className="text-gray-700 hover:text-blue-500 font-semibold">
                                    Submit
                                </Link>
                            </li>
                            <li>
                                <Link to="/execute" className="text-gray-700 hover:text-blue-500 font-semibold">
                                    Execute
                                </Link>
                            </li>
                            <li>
                                <Link to="/schedule" className="text-gray-700 hover:text-blue-500 font-semibold">
                                    Schedule
                                </Link>
                            </li>
                        </ul>
                    </div>
                </nav>
                <div className="pt-16 container mx-auto px-4">
                    <Switch>
                        <Route exact path="/" component={Home} />
                        <Route path="/submit" component={SubmitScript} />
                        <Route path="/execute" component={ExecuteScript} />
                        <Route path="/schedule" component={ScheduleScript} />
                    </Switch>
                </div>
            </div>
        </Router>
    );
}

export default App;