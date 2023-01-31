declare namespace SimpleTimeline {
  interface Entry {
    title: string;
    /**
     * Time is in unix second time
     */
    time: number;
    type: string;
  }
  interface Props {
    list: Entry[];
    bulletSize?: number;
  }
}
