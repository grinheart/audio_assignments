import './App.css';
import Reg from './user/reg';
import Auth from './user/auth';
import Logout from './user/logout';
import CreateTask from './admin/createtask';
import Student from './admin/student';
import Tasks from './admin/tasks';
import StudentsTasks from './student/tasks'
import StudentsTask from './student/task'
import Students from './admin/students';
import AdminAuth from './admin/auth';

import {
  BrowserRouter as Router,
  Switch,
  Route,
} from "react-router-dom";

function App() {
  return (
    <Router>
    <div>
      <Switch>
        <Route exact path="/">
          <div></div>
        </Route>
        <Route path="/reg">
          <Reg />
        </Route>
        <Route path="/auth">
          <Auth />
        </Route>
        <Route path="/logout">
          <Logout />
        </Route>
        <Route exact path='/admin/auth'>
          <AdminAuth />
        </Route>
        <Route exact path='/admin/create'>
          <CreateTask />
        </Route>
        <Route exact path="/admin/student/:id">
          <Student />
        </Route>
        <Route exact path="/admin/students">
          <Students />
        </Route>
        <Route exact path="/admin/tasks">
          <Tasks />
        </Route>
        <Route exact path="/tasks">
          <StudentsTasks />
        </Route>
        <Route exact path="/task/:id">
          <StudentsTask />
        </Route>
      </Switch>
    </div>
  </Router>
  );
}

export default App;
