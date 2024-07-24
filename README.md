# aitime

## aitime is a tool that helps with tracking episodes of TV Shows with the help of AI.
### I use this tool with Siri and iOS Shortcuts because I'm too lazy to open the app and look for myself, so I ask the AI instead.
### The app that is used to track episodes is [TV Time](https://tvtime.com)
### The AI Model that extracts parameters from the given sentence is [codellama](https://ollama.com/library/codellama)

<hr>

### Abilities:
- Add specified episode of a show to a list. "Mark season five episode seven of Game of Thrones as watched" Outputs: ADD_TO_WATCHED:Game of Thrones:5:7
- Remove specified episode of a show to a list. "Remove season five episode seven of Game of Thrones from watched list" Outputs: REMOVE_FROM_WATCHED:Game of Thrones:5:7
- Informs what episode you should watch next for a specified TV Show. "What episode I'm on Game of Thrones Currently?" Outputs: WHERE_WAS_I:Game of Thrones