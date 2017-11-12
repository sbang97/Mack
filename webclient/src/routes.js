import React from "react";
import { Route, IndexRoute } from "react-router";

import Layout from "./components/Layout";
import Home from "./components/Home";
import SignUp from "./components/SignUp";
import NotFound from "./components/NotFound";
import Common from "./components/Common";
import Profile from "./components/Profile";
import Channel from "./components/Channel";

const Routes = (
  <Route path="/" component={Layout}>
    <IndexRoute component={Home} />
    <Route path="/signup" component={SignUp} />
    <Route path="/common" component={Common} />
    <Route path="/profile" component={Profile} />
    <Route path="/channel" component={Channel}/>
    <Route path="*" component={NotFound} />
  </Route>
);

export default Routes;