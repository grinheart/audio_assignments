import './App.css';
import Reg from './user/reg';
import Auth from './user/auth';
import Logout from './user/logout';
import CreateTask from './admin/createtask';
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
        <Route path='/admin/create'>
          <CreateTask />
        </Route>
      </Switch>
    </div>
  </Router>
  );
}

export default App;
