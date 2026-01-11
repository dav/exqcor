import SwiftUI
import ExqcorCore

struct ContentView: View {
    @StateObject private var api = APIClient(baseURL: "http://localhost:3000")
    @State private var baseURL = "http://localhost:3000"
    @State private var selectedScriptId: Int?

    var body: some View {
        VStack(alignment: .leading, spacing: 16) {
            Text("Writer Console").font(.largeTitle)

            VStack(alignment: .leading) {
                Text("Server URL")
                TextField("http://localhost:3000", text: $baseURL)
                    .textFieldStyle(.roundedBorder)
            }

            Button("Load Scripts") {
                api.baseURL = baseURL
                api.loadScripts()
            }

            if let error = api.errorMessage {
                Text(error).foregroundColor(.red)
            }

            HStack(alignment: .top, spacing: 24) {
                VStack(alignment: .leading) {
                    Text("Scripts")
                    List(api.scripts) { script in
                        Button(script.title) {
                            selectedScriptId = script.id
                            api.loadSections(scriptId: script.id)
                        }
                    }
                    .frame(minWidth: 240, minHeight: 200)
                }

                VStack(alignment: .leading) {
                    Text("Sections")
                    List(api.sections) { section in
                        VStack(alignment: .leading) {
                            Text(section.name).font(.headline)
                            if let description = section.description {
                                Text(description).font(.caption)
                            }
                        }
                    }
                    .frame(minWidth: 240, minHeight: 200)
                }
            }

            Spacer()
        }
        .padding()
    }
}
