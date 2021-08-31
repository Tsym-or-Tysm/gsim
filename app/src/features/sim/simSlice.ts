// import { Intent, Position, Toaster } from "@blueprintjs/core";
import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { sendMessage } from "app/appSlice";
import { AppThunk } from "app/store";
import { setActiveName, setLogs, setNames } from "features/debug/debugSlice";
import { setResultData } from "features/results/resultsSlice";

export function saveConfig(path: string, config: string): AppThunk {
  return function (dispatch) {
    const cb = (resp: any) => {
      //check resp code
      if (resp.status !== 200) {
        //do something here
        console.log("Error from server: ", resp.payload);
        return;
      }
      //update
      dispatch(setHasChange(false));
      console.log(resp.data);
    };

    dispatch(
      sendMessage(
        "file",
        "save/file",
        JSON.stringify({
          path: path,
          data: config,
        }),
        cb
      )
    );
  };
}

export function runSim(config: simConfig): AppThunk {
  return function (dispatch, getState) {
    const cb = (resp: any) => {
      dispatch(setLoading(false));
      //check resp code
      if (resp.status !== 200) {
        //do something here
        console.log("Error from server: ", resp.payload);
        // Toaster.create({ position: Position.BOTTOM }).show({
        //   message: "Error running sim: " + resp.payload,
        //   intent: Intent.DANGER,
        // });
        dispatch(setMessage("Sim encountered error: " + resp.payload));
        dispatch(setHasErr(true));

        return;
      }
      //update
      console.log("sim/run received response");
      var data = JSON.parse(resp.payload);
      console.log(data);
      // dispatch(setResultData(data.summary));

      dispatch(setResultData(data));

      if (data.debug) {
        dispatch(setLogs(data.debug));
        dispatch(setNames(data.char_names));
      }

      dispatch(setMessage("Simulation finished. check results"));

      // Toaster.create({ position: Position.BOTTOM }).show({
      //   message: "Simulation finished. check results",
      //   intent: Intent.SUCCESS,
      // });
    };
    dispatch(setLoading(true));
    dispatch(setResultData(null));
    dispatch(setLogs(""));
    dispatch(setNames([]));
    dispatch(setMessage(""));
    dispatch(setHasErr(false));

    //find out who the active is
    const found = config.config.match(/active\+=(\w+);/);
    if (found) {
      dispatch(setActiveName(found[1]));
    }

    dispatch(sendMessage("run", "", JSON.stringify(config), cb));
  };
}

export interface simConfig {
  config: string;
  options: {
    log_details: boolean;
    debug: boolean;
    iter: number;
    workers: number;
    duration: number;
  };
}

export interface ICharacter {
  key: string;
  ascension: number; // 0 to 6
  level: number;
  constellation: number; // 0 to 6
  talentLevelKeys: {
    auto: number;
    skill: number;
    burst: number;
  };
  weapon: IWeapon;
  artifacts: IArtifact; // not used in GO
  stats: number[]; //base stats, should be pulled from api
}

export interface IWeapon {
  key: string;
  level: number; //1 to 90 inclusive
  refinementIndex: number; //0 to 4 inclusive
  stats: number[]; //weapon stats; should be pulled from api
}

export interface IArtifact {
  setKey: string; //set key
  slotKey: string;
  stats: number[]; //artifact stats
  //not used fields
}

interface SimState {
  isLoading: boolean;
  config: string;
  hasChange: boolean;
  msg: string;
  hasErr: boolean;
  characters: ICharacter[];
}
const initialState: SimState = {
  isLoading: false,
  config: "",
  hasChange: false,
  msg: "",
  hasErr: false,
  characters: [],
};

export const simSlice = createSlice({
  name: "sim",
  initialState,
  reducers: {
    setConfig: (state, action: PayloadAction<string>) => {
      state.config = action.payload;
      localStorage.setItem("sim-config", action.payload);
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.isLoading = action.payload;
    },
    setHasChange: (state, action: PayloadAction<boolean>) => {
      state.hasChange = action.payload;
    },
    setMessage: (state, action: PayloadAction<string>) => {
      state.msg = action.payload;
    },
    setHasErr: (state, action: PayloadAction<boolean>) => {
      state.hasErr = action.payload;
    },
  },
});

export const { setConfig, setLoading, setHasChange, setMessage, setHasErr } =
  simSlice.actions;
export default simSlice.reducer;
