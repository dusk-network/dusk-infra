import { useEffect, useRef } from "react";
import { connect } from "react-redux";
import {
  connected,
  connectionError,
  connectTo,
  disconnected,
  updateBlockTimeRead,
  updateCPURead,
  updateDiskRead,
  updateLastBlockInfo,
  updateTxNr,
  updateLogRead,
  updateMemoryRead,
  updateNetRead,
  updateThread,
  updateWarningList,
} from "../../redux/actions";
//
const updateMetrics = {
  cpu: updateCPURead,
  mem: updateMemoryRead,
  latency: updateNetRead,
  disk: updateDiskRead,
};

class DuskSocket {
  constructor(dispatch) {
    this.dispatch = dispatch;
    this.ws = null;
  }

  open(hostname, port) {
    const { dispatch } = this;
    const host = hostname + (port ? ":" + port : "");

    this.ws = new WebSocket(`ws://${host}/stats`);

    ["open", "close", "error", "message"].forEach(type =>
      this.ws.addEventListener(type, this)
    );
  }

  handleEvent(event) {
    const listener = this[`on${event.type}`];
    listener && listener.call(this, event);
  }

  onmessage({ data }) {
    const { dispatch } = this;

    const payload = JSON.parse(data);
    const { metric, text, slice, data: packet, timestamp } = payload;

    switch (metric) {
      case "cpu":
      case "mem":
      case "latency":
      case "disk":
				slice.map(({value, timestamp}) => dispatch(updateMetrics[metric](+value, timestamp)));
        break;

      case "log":
        const { code, level } = packet;
        if (code && code === "round") {
          const { round, blockHash, blockTime, txs } = packet;
          const block = {
            height: round,
            hash: blockHash,
            timestamp,
          };
          dispatch(updateLastBlockInfo(block));
          dispatch(updateBlockTimeRead(blockTime, timestamp));
          dispatch(updateTxNr(txs, timestamp));
          break;
        }
        if (code === "goroutine") {
          const { nr } = packet;

          dispatch(updateThread(nr, timestamp));
          break;
        }

        if (level) {
          const { time = timestamp } = packet;
          dispatch(updateWarningList(packet, time));
          break;
        }
        break;

      case "status":
      	const { warnings=[], round, blockHash, blockTimes=[], txs=[], threads=[] } = packet
      	const block = {
					height: round,
					hash: blockHash,
					timestamp
				};
				dispatch(updateLastBlockInfo(block));
				blockTimes.map(({ value: blockTime, timestamp: stamp }) => dispatch(updateBlockTimeRead(+blockTime, stamp)))
				threads.map(({ value: nr, timestamp: stamp }) => dispatch(updateThread(+nr, stamp)))
				txs.map(({ value: txs, timestamp: stamp }) => dispatch(updateTxNr(+txs, stamp)))
				warnings.map(({ timestamp: stamp, ...others }) => dispatch(updateWarningList(others, timestamp)))

				break;

      case "tail":
        dispatch(updateLogRead(text, timestamp));
        break;

      default:
        console.log("INVALID METRIC SENT");
    }
  }

  onopen() {
    this.dispatch(connected());
  }

  onclose() {
    this.dispatch(disconnected());
  }

  onerror() {
    this.dispatch(connectionError());
  }

  close() {
    this.ws.close();
  }
}

const Socket = ({ connectTo, status, hostname, port, dispatch }) => {
  const ws = useRef();

  useEffect(() => {
    const [hostname, port] = (
      process.env.REACT_APP_HOST_WS || window.location.host
    ).split(":");

    connectTo(hostname, port);

    return () => {
      console.log("cleanup");
    };
  }, []);

  useEffect(() => {
    switch (status) {
      case "connecting":
        ws.current = new DuskSocket(dispatch);
        ws.current.open(hostname, port);
    }
  }, [status]);

  return null;
};

const mapStateToProps = ({ status, hostname, port }) => ({
  status,
  hostname,
  port,
});

const mapDispatchToProps = dispatch => ({
  connectTo: (hostname, port) => dispatch(connectTo(hostname, port)),
  dispatch,
});

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Socket);
