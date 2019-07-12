import AppBar from "@material-ui/core/AppBar";
import Container from "@material-ui/core/Container";
import CssBaseline from "@material-ui/core/CssBaseline";
import Grid from "@material-ui/core/Grid";
import Paper from "@material-ui/core/Paper";
import { makeStyles } from "@material-ui/core/styles";
import Toolbar from "@material-ui/core/Toolbar";
import Typography from "@material-ui/core/Typography";
import clsx from "clsx";
import React from "react";
import { connect } from "react-redux";
import logo from "../../d.svg";
import {
  getCPUMetrics,
  getCurrentBlockInfo,
  getDiskMetrics,
  getLogMetrics,
  getMemoryMetrics,
  getNetMetrics,
  getTimeMetrics,
  getWarnings,
} from "../../redux/selectors";
import BlockCreated from "./BlockCreated";
import BlockHeight from "./BlockHeight";
import BlockTimeChart from "./BlockTimeChart";
import BlockTransactionChart from "./BlockTransactionChart";
import CPUChart from "./CPUChart";
import DiskChart from "./DiskChart";
import LogFile from "./LogFile";
import MemChart from "./MemChart";
import ThreadChart from "./ThreadChart";
import NetLatencyChart from "./NetLatencyChart";
import Warnings from "./Warnings";

const logify = (text, { D }) => {
  return text
    .toUpperCase()
    .split(/\b(D)/)
    .map(item =>
      item !== "D" ? item : <img src={logo} alt="D for Dusk" className={D} />
    );
};

function Copyright() {
  return (
    <Typography variant="body2" color="textSecondary" align="center">
      {"Copyright © 2019 Dusk Network B.V. All Rights Reserved."}
    </Typography>
  );
}

const drawerWidth = 240;

const useStyles = makeStyles(theme => ({
  root: {
    display: "flex",
  },
  toolbar: {
    paddingRight: 24, // keep right padding when drawer closed
  },
  toolbarIcon: {
    display: "flex",
    alignItems: "center",
    justifyContent: "flex-end",
    padding: "0 8px",
    ...theme.mixins.toolbar,
  },
  appBar: {
    zIndex: theme.zIndex.drawer + 1,
    transition: theme.transitions.create(["width", "margin"], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
  },
  appBarShift: {
    marginLeft: drawerWidth,
    width: `calc(100% - ${drawerWidth}px)`,
    transition: theme.transitions.create(["width", "margin"], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
  },
  menuButton: {
    marginRight: 36,
  },
  menuButtonHidden: {
    display: "none",
  },
  D: {
    height: ".95em",
    marginRight: ".4em",
    paddingTop: 1,
    alignSelf: "center",
  },
  title: {
    fontFamily: "Lato",
    textTransform: "uppercase",
    letterSpacing: ".5em",
    flexGrow: 1,
    display: "flex",
    whiteSpace: "pre-wrap",
  },
  drawerPaper: {
    position: "relative",
    whiteSpace: "nowrap",
    width: drawerWidth,
    transition: theme.transitions.create("width", {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
  },
  drawerPaperClose: {
    overflowX: "hidden",
    transition: theme.transitions.create("width", {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
    width: theme.spacing(7),
    [theme.breakpoints.up("sm")]: {
      width: theme.spacing(9),
    },
  },
  appBarSpacer: theme.mixins.toolbar,
  content: {
    flexGrow: 1,
    height: "100vh",
    overflow: "auto",
  },
  container: {
    paddingTop: theme.spacing(4),
    paddingBottom: theme.spacing(4),
  },
  paper: {
    padding: theme.spacing(2),
    display: "flex",
    overflow: "auto",
    flexDirection: "column",
  },
  fixedHeight: {
    height: 240,
  },
  noScrolling: {
    overflow: "hidden",
  },
}));

function Dashboard({
  hostname,
  items,
  lastBlock,
  blockTime,
  net,
  disk,
  memory,
  log,
  score,
  cpu,
  warnings,
}) {
  const classes = useStyles();
  const [open, setOpen] = React.useState(false);
  const handleDrawerOpen = () => {
    setOpen(true);
  };
  const handleDrawerClose = () => {
    setOpen(false);
  };
  const fixedHeightPaper = clsx(classes.paper, classes.fixedHeight);
  const fixedHeightPaperNoScrollig = clsx(
    fixedHeightPaper,
    classes.noScrolling
  );

  return (
    <div className={classes.root}>
      <CssBaseline />
      <AppBar
        position="absolute"
        className={clsx(classes.appBar, open && classes.appBarShift)}
      >
        <Toolbar className={classes.toolbar}>
          <Typography
            component="h1"
            variant="h6"
            color="inherit"
            noWrap
            className={classes.title}
          >
            {logify("Duskboard", classes)}
          </Typography>
          <Typography component="h1" variant="h6" color="inherit" noWrap>
            {hostname}
          </Typography>
        </Toolbar>
      </AppBar>

      <main className={classes.content}>
        <div className={classes.appBarSpacer} />
        <Container maxWidth="lg" className={classes.container}>
          <Grid container spacing={3}>
            <Grid item xs={12} sm={7}>
              <Paper className={classes.paper}>
                <BlockHeight height={lastBlock.height} hash={lastBlock.hash} />
              </Paper>
            </Grid>
            <Grid item xs={12} sm={5}>
              <Paper className={classes.paper}>
                <BlockCreated timestamp={lastBlock.timestamp} />
              </Paper>
            </Grid>
            <Grid item xs={12} sm={6}>
              <Paper className={fixedHeightPaper}>
                <BlockTimeChart data={blockTime} />
              </Paper>
            </Grid>
            <Grid item xs={12} sm={6}>
              <Paper className={fixedHeightPaper}>
                <BlockTransactionChart data={blockTime} />
              </Paper>
            </Grid>
            <Grid item xs={12} sm={9}>
              <Paper className={fixedHeightPaperNoScrollig}>
                <ThreadChart data={memory} />
              </Paper>
            </Grid>
            <Grid item xs={12} sm={3}>
              <Paper className={fixedHeightPaperNoScrollig}>
                <DiskChart data={disk} />
              </Paper>
            </Grid>
            <Grid item xs={12} sm={4}>
              <Paper className={fixedHeightPaperNoScrollig}>
                <CPUChart data={cpu} />
              </Paper>
            </Grid>

            <Grid item xs={12} sm={4}>
              <Paper className={fixedHeightPaperNoScrollig}>
                <MemChart data={memory} />
              </Paper>
            </Grid>
            <Grid item xs={12} sm={4}>
              <Paper className={fixedHeightPaperNoScrollig}>
                <NetLatencyChart data={net} />
              </Paper>
            </Grid>
            <Grid item xs={12}>
              {/* Error lists to be added to a board to be less ephemeral */}
              <Warnings items={warnings} />
            </Grid>
            <Grid item xs={12}>
              <LogFile items={log} />
            </Grid>
          </Grid>
        </Container>
        <Copyright />
      </main>
    </div>
  );
}

const mapStateToProps = state => ({
  hostname: state.hostname,
  lastBlock: getCurrentBlockInfo(state),
  // items: lastNodeUpdateSelector(state),
  // locations: getNodeLocations(state),
  // score: getHighestScore(state),
  blockTime: getTimeMetrics(state),
  cpu: getCPUMetrics(state),
  log: getLogMetrics(state),
  net: getNetMetrics(state),
  disk: getDiskMetrics(state),
  memory: getMemoryMetrics(state),
  warnings: getWarnings(state),
});

export default connect(mapStateToProps)(Dashboard);
