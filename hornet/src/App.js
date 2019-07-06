import React from "react";
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
import Dashboard from "./components/dashboard/Dashboard";
import Socket from "./components/socket/Socket";
export default function App({ listenForUpdates }) {
  return (
    <>
      <Socket />
      <Dashboard />
    </>
  );
}
