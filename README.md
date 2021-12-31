# Modeling Epidemics Using The SEIR Model
### **PROBELM AND MOTIVATION**
Epidemics are tricky and can spread quickly if not acted upon right away. In the era of COVID, modeling epidemics to track the way they move has been more important than ever to make better decisions on how to act in response to a spread. Through modeling different scenarios such as hotspots, quarantining, etc., we can better prepare for many possibilities of what could happen when an infectious diease spread. Additionally, we can better understand how preventative measures such as handwashing, quarantining, etc. can help slow down the spread.

### **DESCRIPTION**

This project allows the user to model a spread of an epidemic by manipulating different parameters based off the SEIR model. From mimicking different communities such as neighborhoods quarantining to countries closing off borders, this program gives the user a lot of flexibility in choosing their community and parameters to model. 

#### **FEATURES**
* Generates an gif animiating the epidemic spread
* Generates the first image, depicting the start of the epidemic
* Generates a scattorplot following the SEIR model

#### **BENCHMARK**

If my simulation generates a similar graph to the SEIR model and can model different scenarios such as qurantine and hotspots, it was successful!

#### **HOW IT WORKS**
The program starts by initializing people on a board, deciding whether they are infectious or not based on the starting infection rate. If it specified to have communities, the board will be split into the number of specified communities. People are randomly placed onto the board and  squares can be specified to have only a certain number of people in them. Squares can be chosen to closed off (meaning no one can enter and the people inside cannot leave) or just qurantining (meaning people can still enter, but the people inside cannot leave). The amount of people leaving and entering a square can be changed to mimick strict vs not so strict border rules. Based on the specified exposure, infection, recover, and lethality rates, the status of people on the board will change based on who they come in contact with. A contact range is specified (aoe), and if someone is within that range, they "interact" with the person in range. The higher the rate, the more likely someone is to be exposed, infected, etc. The board is run through a certain number of generations, or think of this like time. Each board generation is seen animated in a gif showing the overall spread. To generate the scatterplot, the amount of susceptible, exposed, recovered, and infected people are counted and then plot over the number of generations. 

To learn more about SEIR models see:
* Carcione JM, Santos JE, Bagaini C, Ba J. A Simulation of a COVID-19 Epidemic Based on a Deterministic SEIR Model. Front Public Health. 2020;8:230. Published 2020 May 28. doi:10.3389/fpubh.2020.00230

### **ISSUES**
* Movement: Although there are many checks to keep people on the board or in and out of squares, since people are moved randomly, sometimes people pop up in places they should not be. 
* Mixing flags and file: The file parameters should overtake the flag parameters, however some parameters have not been specified for file format yet and could cause potential issue when mixing file and flag parameters.

### **HOW TO USE**
Make sure you have the following installed:
* github.com/Arafatk/glot (use "go get github.com/Arafatk/glot")
*  gnuplot (use "brew install gnuplot" - Mac)

There are two ways to create a simulation.
1. Use flags in terminal to change different parameters for the simluation
2. Use a combination of flags and a txt file to specify the starting board and parameters for the board

The first method is simplier to use as there is no need to create a txt file. However, the first method yields simplier boards and simulations. The second method is a bit more complext to use as you have to create a txt file to specify your board but, we can create different shapes and more interesting scenarios. See "Parameters" for all the flags you can change and "Try Out Different Board" for the specifics on what the txt file should look like. 

### **PARAMETERS**
Use these flags to change the following parameters to simulate an epidemic:

| Flag        | Default Value | Description |
| ----------- | ----------- | ----------- |
numGen | 100| Number of steps to run the community
s| 10| Number of squares or slices to split board
aoe| 1|Radius distance of infection | unit: slices/squares
move| 50| How much should people move?
| width| 1000| Width (and height) of the image to create
sf| 10| a scaling factor for size of people
si| 0.01| a starting % of the population that is infected
|e| 0.235 |likelihood for getting exposed given that are come in contact with someone infected
i |0.157 |likelihood of turning infectious at a time step 
r |0.97 |likelihood of recovering at a time step
minDay |14 |min days infectious |
l |0.027 |likelihood of dying
o| out.png |Name of PNG to output
a | anim | Animated GIF to write
freq| 1| frame writing interval
numPeop | 100 |Number of people on the board"
cBool | false |Create communities (ie. countries), # of communites: slice x slice
numClose | 0 |Number of communities to close borders and quarantine
numQ | 5 | Number of communities to quarantine but aren't closing borders 
f| "" |Name of file to create communities
clus|false |Should there be a cluster on the board?
clusX| 1|What x coordinate should the people cluster at?"
clusY| 1|What y coordinate should the people cluster at?"
clusDay| 5|How many days should a person spend in a hot spot?


#### **TRY OUT DIFFERENT BOARD**
Create a .txt file in the following format:

| Format | Example (3 x 3 square) |
--------------------------------------------------------------------------------------- |  ------------------------ 
row 0: board dimension                                                                  |   3
row 1: square 1(x,y) square 2(x,y) square 3(x,y)                                        |   0 1 2 1 2 2
row 2: square 1(% entry) square 2(% entry) square 3(% entry)                            |   1 0.5 0.01
row 3: square 1(% leave) square 2(% leave) square 3(% leave)                            |   0.3 0.1 0.8
row 4: square 1(initial % infect) square 2(initial % infect) square 3(initial % infect) |   0.05 0.2 0.8
row 5: square 1(# ppl) square 2(# ppl) square 3(# ppl)                                  |   5 6 7
row 6: cluster square 1(x,y) cluster square 2 (x,y)                                     |   0 0 0 2
row 7: cluster square time 1 cluster square time 2                                      |   3 5


##### **NOTES**
* MUST include row 0, row 1, row 2, row 3
* Squares to cluster in CANNOT be the same square that has an entry % of 0
* If row 6 is specified, so MUST be row 7
* Board will always be a square, number indicates # of squares for one side
* Cluster in squares (0,0) and (0,2), with 3 days spent in cluster 1 and 5 days spent in cluster 2

##### **Where Square 1:**
* x, y: 0,1
* % entry: 1
* % leave: 0.3
* initial % infection: 0.05
* Number of people in square: 5



##### **Entry and Leaving Percentages Mean:**
* 0: No one can enter or leave
* 1: Anyone can leave or enter a square

See the .txt files for additional reference
