import React, { useEffect } from "react";
import Container from "@material-ui/core/Container";
import Typography from "@material-ui/core/Typography";
import Box from "@material-ui/core/Box";
import ProTip from "./ProTip";
import Link from "@material-ui/core/Link";
import Dashboard from "./components/dashboard/Dashboard";
import Socket from "./components/socket/Socket";
// import React from "react";
// import { connect } from "react-redux";
// const mapStateToProps = state => {
//   return { articles: state.articles };
// };
// const ConnectedList = ({ articles }) => (
//   <ul className="list-group list-group-flush">
//     {articles.map(el => (
//       <li className="list-group-item" key={el.id}>
//         {el.title}
//       </li>
//     ))}
//   </ul>
// );
// const List = connect(mapStateToProps)(ConnectedList);
// export default List;

import "./chartist.scss";
export default function App({ listenForUpdates }) {
  return (
    <>
      <Socket />
      <Dashboard />
    </>
  );
}
