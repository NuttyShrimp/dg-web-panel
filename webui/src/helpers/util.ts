export const parsePotentialJSON = (str: string | number | Record<string, string>) => {
  try {
    return JSON.parse(String(str));
  } catch (e) {
    return str;
  }
}
