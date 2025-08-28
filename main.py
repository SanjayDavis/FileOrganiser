import os
import json
import shutil

def load_mappings(filename):
    # load the json file to read
    try:
        with open(filename, "r") as f:
            return json.load(f)
    except FileNotFoundError:
        print(f"Error: {filename} not found")
        return {}
    except json.JSONDecodeError as e:
        print(f"Error parsing {filename}: {e}")
        return {}

def list_files(path):
    try:
        return os.listdir(path)
    except OSError as e:
        print(f"Error reading directory {path}: {e}")
        return []

def get_unique_path(dest):
    if not os.path.exists(dest):
        return dest
    
    base, ext = os.path.splitext(dest)
    counter = 1
    new_dest = f"{base}({counter}){ext}"
    while os.path.exists(new_dest):
        counter += 1
        new_dest = f"{base}({counter}){ext}"
    return new_dest

def move_file(src, dest):
    try:
        dest = get_unique_path(dest)
        print(f"Moving file {src} to {dest}")
        shutil.move(src, dest)
    except Exception as e:
        print(f"Error moving {src} to {dest}: {e}")

def main():
    path = os.getcwd()
    print("Current Working Directory is:", path)
    print("Checking for all files in the directory")

    files = list_files(path)
    dir_map = load_mappings("mappings.json")

    for filename in files:
        # Skip config and script file itself
        if filename in ("mappings.json", os.path.basename(__file__)):
            continue

        file_path = os.path.join(path, filename)

        if os.path.isfile(file_path):
            ext = os.path.splitext(filename)[1].lower().strip()
            dir_name = dir_map.get(ext, "others")
            dir_path = os.path.join(path, dir_name)

            os.makedirs(dir_path, exist_ok=True)

            new_path = os.path.join(dir_path, filename)
            move_file(file_path, new_path)

if __name__ == "__main__":
    main()
