package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)


//SmartContract that contain the Contract API , we will build our function on this struct
type SmartContract struct {
	contractapi.Contract
}

//User struct contain the name and surname of a User
type User struct {
	Name 			string `json:"Name"`
	Surname			string `json:"Surname"`
}

//Presence Struct contain all the information that we want to validate in the World State. The ID is used as the key 
type Presence struct {
	ID             string `json:"ID"`
	User 		   User	  `json:"User"`
	NameTag        string `json:"nametag"`
	DiscordID	   string `json:"discordid"`
	TwitchID       string `json:"twitchid"`
	ValidatedAccount	bool `json:"validatedaccount"`
	Condition	   int	  `json:"Condition"`

}

//SetConditionUsed will set the condition to 1 , to notify that the asset has been modifier
func (Presence *Presence) SetConditionUsed(){
	Presence.Condition = 1
}

//InitLedger marshall some stub ID and put them in the world state. The ID is used as the key 
func (sc *SmartContract) InitLedger( ctx contractapi.TransactionContextInterface) error {
	StubID := []Presence{
		{ID: "JD2020230001",User:User{Name:"John",Surname:"Doe"},NameTag:"SuperUser",DiscordID:"DiscordID92",TwitchID:"SKDIAS",ValidatedAccount: true,Condition: 0},
		{ID: "OD2020230001",User:User{Name:"Oliver",Surname:"Tree"},NameTag:"Pwned",DiscordID:"DiscordID03",TwitchID:"SKDSSAS",ValidatedAccount: false,Condition: 0},
	}
	for _,Identity := range StubID {
		IdentityJSON, err := json.Marshal(Identity)

		if err != nil {
			return err
		}

		ctx.GetStub().PutState(Identity.ID,IdentityJSON)

		if err != nil {
			return fmt.Errorf("Failed to set inital value on worldState. Returned error: %s", err)
		}
	}

	return nil 
}

//CreateID function is used to add an ID to the World state , it take all the elements of and ID , Marshall them and put them in the world state
func (sc *SmartContract) CreateID(ctx contractapi.TransactionContextInterface, ID string , surname string , name string , nameTag string , discordID string , twitchID string, ValidatedAccount bool, Condition int) error {
	
	Identity :=Presence{
		ID: ID ,
		User: User{
			Name: name,
			Surname: surname,
		},
		NameTag: nameTag,
		DiscordID: discordID,
		TwitchID: twitchID,
		ValidatedAccount: ValidatedAccount,
		Condition: 0,	
	}
	existingID , err :=ctx.GetStub().GetState(ID)

	if err != nil {
		return fmt.Errorf("Failed to interact with world state , error returned: %s" , err)
		
	}

	if existingID != nil {
		return fmt.Errorf("Cannot create Identity with Given ID , one already present in the world state with value: %s", existingID)
	}


	IdentityJSON, err := json.Marshal(Identity)

	if err != nil {
		return err
	}

	ctx.GetStub().PutState(Identity.ID,IdentityJSON)

	if err != nil {
		return fmt.Errorf("Failed to set ID in World State , error returned: %s", err)
	}
	return nil

}

//GetID function will query the Identity with a given key
func (sc *SmartContract) GetID(ctx contractapi.TransactionContextInterface, ID string) (*Presence,error) {

	identityAsByte , err := ctx.GetStub().GetState(ID)
	if err != nil {
		return nil, fmt.Errorf("Can't interact with world state , returned error : %s" ,err)
	}

	if identityAsByte == nil {
		return nil , fmt.Errorf("Can't retrieve specified ID , value does not exist")
	}

	identity := new(Presence)

	_ = json.Unmarshal(identityAsByte,identity)

	return identity, nil
}

//ChangeNameTag func will change a value in a Identity , it return the error if one occurs
func (sc *SmartContract) ChangeNameTag(ctx contractapi.TransactionContextInterface , ID string, newValue string) error {
	
	identity , err := sc.GetID(ctx,ID)

	if err != nil {
		return fmt.Errorf("Can't interact with world state , returned error : %s" ,err)
	}

	if identity == nil {
		return fmt.Errorf("Can't find any identity with ID %s in the world state. erro: %s",ID,err)
	}

	identity.NameTag = newValue
	identity.SetConditionUsed()

	identityAsByte, _ := json.Marshal(identity)


	return ctx.GetStub().PutState(ID,identityAsByte)

}


func main() {
	chaincode , err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Can't creating chaincode: %s", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode , returned error: %s",err)
	}

}
