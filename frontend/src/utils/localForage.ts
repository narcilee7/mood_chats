import localforage from "localforage";


const getValue = async (key: string): Promise<any> => {
  let value = null;
  try {
    value = await localforage.getItem(key);
  } catch (error) {
    console.error("Error getting value from localForage:", error);
    return null;
  }
  return value ? JSON.parse(value as string) : null;
};


const setValue = async (key: string, value: any): Promise<void> => {
  try {
    await localforage.setItem(key, JSON.stringify(value));
  } catch (error) {
    console.error("Error setting value in localForage:", error);
  }
};

export {
  getValue,
  setValue
}