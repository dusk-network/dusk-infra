import { connect } from "react-redux";
import { listenForUpdates } from "../../redux/actions";
import { useEffect } from "react";

function mapDispatchToProps(dispatch) {
  return {
    listenForUpdates: socket => dispatch(listenForUpdates(socket))
  };
}

const socket = {};

const Socket = ({ listenForUpdates }) => {
  useEffect(() => {
    console.log("mounted");
    listenForUpdates(socket);

    return () => {
      console.log("cleanup");
    };
  });

  return null;
};

export default connect(
  null,
  mapDispatchToProps
)(Socket);
