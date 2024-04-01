import { ChangeEvent, MouseEvent, useState } from "react";
import "./App.css";

function App() {
  const [file, setFile] = useState<File>();

  function handleChange(event: ChangeEvent<HTMLInputElement>) {
    const target = event.target;
    if (target && target.files!.length > 0) {
      setFile(target.files![0]);
    }
  }

  const handleResponse = () => {};
  const handleError = () => {};

  const handleSubmit = (e: MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    if (file) {
      const url = "/api/file/upload";
      const formData = new FormData();
      formData.append("file", file);
      formData.append("fileName", file.name);
      const headers = {
        "content-type": "multipart/form-data",
      };
      fetch(url, { headers, method: "POST", body: formData })
        .then(handleResponse)
        .catch(handleError);
    }
  };

  return (
    <div className="App">
      <form>
        <input
          type="file"
          name="img_upload"
          id="img_upload"
          onChange={handleChange}
        />
        <button onClick={handleSubmit}>Upload</button>
      </form>
    </div>
  );
}

export default App;
