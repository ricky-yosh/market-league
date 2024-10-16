from roadmapper.roadmap import Roadmap
from roadmapper.timelinemode import TimelineMode

# Initialize the Roadmap
roadmap = Roadmap(2000, 1000, colour_theme="BLUEMOUNTAIN")
roadmap.set_title("MarketLeague Roadmap v1.0")
roadmap.set_timeline(TimelineMode.WEEKLY, start="2024-08-26", number_of_items=15)

# Group 1: Documentation
group = roadmap.add_group("Documentation", fill_colour="#FFC000", font_colour="black")

task = group.add_task("Problem Research", "2024-08-26", "2024-09-05")

task = group.add_task("Planning Ideas w/ Professor", "2024-08-29", "2024-09-08")
task.add_milestone("Idea Finalized", "2024-09-08")

task = group.add_task("Use Case, Sequence, Architecture Diagrams", "2024-09-09", "2024-09-22")
task = group.add_task("Dataflow, Class, Database Diagrams", "2024-09-19", "2024-09-30")
task.add_milestone("Initial Diagrams Created", "2024-09-30")

task = group.add_task("Roadmap Creation", "2024-10-01", "2024-10-13")
task = group.add_task("Iterate on Diagrams", "2024-10-10", "2024-12-04")

# Group 2: Technical Development
group = roadmap.add_group("Technical Development", fill_colour="#70AD47", font_colour="black")

task = group.add_task("Technology and Platform Research", "2024-08-26", "2024-09-15")
task.add_milestone("Technologies Selected", "2024-09-15")

task = group.add_task("Integration Testing (Angular, Gin, Postgres)", "2024-09-16", "2024-10-06")
task.add_milestone("Basic Button Integration Success", "2024-10-06")

task = group.add_task("Stock API Testing", "2024-10-07", "2024-10-20")
task = group.add_task("Login/Signup", "2024-10-07", "2024-10-20")
task = group.add_task("User Class", "2024-10-07", "2024-10-21")

# Future Tasks
task = group.add_task("Stock Objects", "2024-10-10", "2024-10-26")
task = group.add_task("Portfolio Functionality", "2024-10-10", "2024-10-27")
task = group.add_task("League Creation", "2024-10-11", "2024-10-28")
task = group.add_task("Trade Functionality", "2024-10-11", "2024-10-28")
task = group.add_task("Trade History", "2024-10-11", "2024-10-28")

task = group.add_task("Scoring System Implementation", "2024-10-28", "2024-11-14")
task = group.add_task("Leaderboard", "2024-10-30", "2024-11-17")
task = group.add_task("Trade Fairness Algorithm", "2024-10-31", "2024-11-18")
task = group.add_task("Stock Charts Data Visualization", "2024-11-16", "2024-12-02")

# Group 3: Presentations
group = roadmap.add_group("Presentations", fill_colour="#ED7D31", font_colour="black")

task = group.add_task("Idea Presentation Work", "2024-08-26", "2024-09-11")
task.add_milestone("Idea Proposal", "2024-09-11")

task = group.add_task("Diagrams Presentation Preparation", "2024-09-14", "2024-10-02")
task.add_milestone("Diagrams & Documents Presentation", "2024-10-02")

task = group.add_task("Prototype Slides and Demo Work", "2024-10-05", "2024-10-30")
task.add_milestone("Prototype Demo", "2024-10-30")

task = group.add_task("Alpha Demonstration Preparation", "2024-11-02", "2024-12-04")
task.add_milestone("Alpha Demo", "2024-12-04")

# Generate Roadmap
roadmap.draw()
roadmap.save("marketleague-roadmap.png")
