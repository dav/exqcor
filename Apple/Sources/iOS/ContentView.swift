import SwiftUI
import ExqcorCore

struct ContentView: View {
    @StateObject private var api = APIClient(baseURL: "http://localhost:3000")
    @State private var baseURL = "http://localhost:3000"
    @State private var selectedScriptId: Int?

    var body: some View {
        NavigationView {
            VStack(alignment: .leading, spacing: 16) {
                Text("Writer Console").font(.title2)

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

                List {
                    Section("Scripts") {
                        ForEach(api.scripts) { script in
                            Button(script.title) {
                                selectedScriptId = script.id
                                api.loadSections(scriptId: script.id)
                            }
                        }
                    }

                    Section("Sections") {
                        ForEach(api.sections) { section in
                            VStack(alignment: .leading) {
                                Text(section.name).font(.headline)
                                if let description = section.description {
                                    Text(description).font(.caption)
                                }
                            }
                        }
                    }
                }
            }
            .padding()
            .navigationTitle("Exqcor")
        }
    }
}
